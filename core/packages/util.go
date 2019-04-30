package packages

import (
  "archive/tar"
  "compress/gzip"
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "os"
  "path/filepath"
  "runtime/debug"
  "strconv"
  "strings"
  "time"
)

/*
  Downloads a file with through an HTTP request.
*/
func downloadFileWithHttp(filename string, source string) {
  success := false
  // We've occasionally seen issues when downloading files, for
  // example "connection reset by peer" which can happen when file
  // servers are overloaded with download requests from this client
  // code. It's not always clear what the root cause is, so we're
  // extra conservative with retrying in most cases in the following
  // code, even though some errors will be permanent rather than
  // transient. This code is a bit messy and brittle and can
  // definitely be improved upon in the future as we learn more about
  // the common error cases.
  for attempts := 1; attempts <= 3; attempts++ {
    fmt.Println("\nDownloading file from:", source, "Attempt number:", attempts)
    // Make an http request to get the package file
    request, err := http.NewRequest("GET", source, nil)
    if err != nil {
      fmt.Println("Error with GET request:", err, "Retrying...")
      time.Sleep(3 * time.Second)
      continue
    }

    resp, err := http.DefaultClient.Do(request)
    if err != nil {
      fmt.Println("Error with sending http request:", err, "Retrying...")
      time.Sleep(3 * time.Second)
      continue
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
      fmt.Println("Http response status:", resp.Status)
    }

    // Start a download progress printer coroutine, channel used to signal it to end
    signal := make(chan bool)
    fileSize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
    go printDownloadProgress(filename, fileSize, signal)

    // Write the contents of body (the package) to a file
    err = createFile(filename, 0644, resp.Body)
    // Done with download
    signal <- true
    if err != nil {
      fmt.Println("Error with saving downloaded file:", err, "Retrying...")
      time.Sleep(3 * time.Second)
      continue
    }
    success = true
    break
  }
  if !success {
    handleErr(fmt.Errorf("Maximum retry limit reached. Failed to download file: %s. Exiting.", filename))
  }
}

/*
  Prints download progress messages for a given file.
*/
func printDownloadProgress(filename string, fileSize int, signal chan bool) {
  // Continues to print until item is added to signal channel
  for {
    // Either end or keep printing progress messages
    select {
    case <-signal:
      return
    default:
      // Get the current size of the target file and print a percent based on it and the expected file size
      file, err := os.Open(filename)
      if err != nil {
        continue
      }
      fileInfo, err := file.Stat()
      if err != nil {
        handleErr(err)
      }
      size := fileInfo.Size()

      percent := float64(size) / float64(fileSize) * 100
      fmt.Printf(getProgressText(percent), percent)
    }
    time.Sleep(time.Second)
  }
}

func getProgressText(percent float64) string {
  progressString := "\r"
  maxNumBars := 50
  numBars := int(percent) / (100 / maxNumBars)
  for i := 1; i <= maxNumBars; i++ {
    switch {
    case i <= numBars:
      progressString += "#"
    default:
      progressString += " "
    }
  }
  progressString += " (%.2f %%)"
  return progressString
}

func displayEulaAndAskForAcknowledgement(eula_text string, eula_file_path_to_display string) bool {
  fmt.Printf("\n\nPlease read the following agreement. If you'd like to re-read it again")
  fmt.Printf("\nin the future, the below text will be at:\n%s", eula_file_path_to_display)
  fmt.Printf("\n\n%s", eula_text)
  fmt.Printf("\n\nDo you agree to the above terms? [yes/no]: ")
  response := ""
  _, err := fmt.Scanln(&response)
  fmt.Printf("\nYour response: %s\n", response)
  // The following loop is kinda hacky. This is because golang handles
  // spaces in input text a bit weird, delimiting on each one. So if
  // the user typed e.g. "foo bar<ENTER>", the result is that only
  // "foo" was read into &response, and "bar" is still hanging around
  // on stdin. The Scanln function sets a particular err value in this
  // case, because it always expects a new line character after
  // reading into &response. Here we condition on the err to check for
  // that case, in order to clear any extra user input such as the
  // example "bar", to make sure it doesn't unintentionally get read
  // into any future prompts for user input. It's possible there's a
  // better way to do this, but I spent awhile digging around online
  // and seems lots of folks have such issues with golang user input
  // that contains spaces, and there isn't a great solution. The often
  // cited bufio.NewReader method is pretty inconvenient here, because
  // as far as I can tell, it usually reads in more data than you'd
  // want to process, for example if a long-ish input string was piped
  // into mcli rather than typed interactively, we usually would only
  // want to process up to the first new line in this invocation of
  // getEulaAcknowledgementAndExtractTarball, leaving the rest still
  // on stdin. The bufio way would read all the data, leaving nothing
  // left.
  for err != nil && err.Error() == "expected newline" {
    var s string
    _, err = fmt.Scanln(&s)
    // We just discard the input that was read. We'll keep doing that
    // until we hit a new line.
  }
  response = strings.ToLower(response)
  return response == "y" || response == "yes"
}

func CheckOrAskForMcliEulaAcknowledgement(baseDir string) bool {
  acknowledgement_file_path := filepath.Join(baseDir, "etc/mcli/eula_acknowledgement.txt")
  expected_ack_string := "I agree to the MCLI end user license agreement.\n"
  ack_bytes, err := ioutil.ReadFile(acknowledgement_file_path)
  if err == nil && string(ack_bytes) == expected_ack_string {
    return true
  }

  eula_path := filepath.Join(baseDir, "EULA.txt")
  eula_bytes, err := ioutil.ReadFile(eula_path)
  if err != nil {
    fmt.Printf("\nFailed to find EULA.txt file for MCLI. Please reinstall the MCLI.\n")
    return false
  }
  fmt.Printf("\nHi! It looks like you're either running this version of the CLI\nsoftware for the first time or haven't yet acknowledged the license.\n")
  if !displayEulaAndAskForAcknowledgement(string(eula_bytes), eula_path) {
    fmt.Printf("\nExiting since you did not agree.\n")
    return false
  }

  err = os.MkdirAll(filepath.Dir(acknowledgement_file_path), 0755)
  if err == nil {
    err = ioutil.WriteFile(acknowledgement_file_path, []byte(expected_ack_string), 0644)
  }
  if err != nil {
    fmt.Printf("\nFailed to record eula acknowledgement:\n%v\n", err)
    fmt.Printf("\nThis means you'll be shown the EULA again the next time you run the program.")
    fmt.Printf("\nPerhaps check the write permissions on your MCLI installation directory to fix this for good.\n")
  }
  return true
}

/*
  Looks for a *eula.txt file in a tarball, prints it and gets user
  acknowledgement, then extracts everything in the tarball file to a
  target directory.
*/
func getEulaAcknowledgementAndExtractTarball(tarballName string, targetDir string) {
  // It seems tar format only allows sequential access to data, so we
  // end up reading the tarball twice, once to look for a *eula.txt
  // file and once again to actually extract the data. This is because
  // we want to be careful not to extract the data if the user doesn't
  // acknowledge the eula.
  eula_filename, eula_text, err := extractEula(tarballName)
  if err != nil {
    fmt.Println("FATAL ERROR: Failed to extract a EULA from downloaded package. Exiting.")
    handleErr(err)
  }
  eula_path := filepath.Join(targetDir, eula_filename)
  if displayEulaAndAskForAcknowledgement(eula_text, eula_path) {
    fmt.Printf("\nInstalling package.\n")
    extractTarball(tarballName, targetDir)
  } else {
    fmt.Printf("\nSkipping installation of this package since you did not agree.\n")
  }
}

func extractEula(tarballName string) (string, string, error) {
  file, err := os.Open(tarballName)
  if err != nil {
    handleErr(err)
  }

  // file is tar'ed then gzipped, we need to do the reverse
  // gzip reader
  gzReader, err := gzip.NewReader(file)
  if err != nil {
    handleErr(err)
  }
  defer gzReader.Close()
  // tar reader
  tarReader := tar.NewReader(gzReader)

  for {
    // tar's reader returns next entry in tar, or EOF if done
    header, err := tarReader.Next()
    if err != nil {
      return "", "", err
    }

    lowered_name := strings.ToLower(header.Name)
    if strings.HasSuffix(lowered_name, "eula.txt") &&
      header.FileInfo().Mode().IsRegular() {
      content, err := ioutil.ReadAll(tarReader)
      if err != nil {
        handleErr(err)
      }
      return header.Name, string(content), nil
    }
  }
  return "", "", err
}

func extractTarball(tarballName string, targetDir string) error {
  fmt.Println("\nUntaring file: ", tarballName)

  file, err := os.Open(tarballName)
  if err != nil {
    handleErr(err)
  }

  // file is tar'ed then gzipped, we need to do the reverse
  // gzip reader
  gzReader, err := gzip.NewReader(file)
  if err != nil {
    handleErr(err)
  }
  defer gzReader.Close()
  // tar reader
  tarReader := tar.NewReader(gzReader)

  for {
    // tar's reader returns next entry in tar, or EOF if done
    header, err := tarReader.Next()
    switch {
    case err != nil:
      return err
    case err == io.EOF:
      return nil
    }

    // create dir or create file
    targetFile := filepath.Join(targetDir, header.Name)
    switch {
    case header.FileInfo().Mode().IsDir():
      createDir(targetFile)
    case header.FileInfo().Mode().IsRegular():
      createFile(targetFile, os.FileMode(header.Mode), tarReader)
    case header.FileInfo().Mode()&os.ModeSymlink != 0:
      createSymlink(header.Linkname, targetFile)
    }
  }
}

/*
  Creates a directory at target path
*/
func createDir(path string) error {
  if _, err := os.Stat(path); err != nil {
    if err := os.MkdirAll(path, 0755); err != nil {
      return err
    }
  }
  fmt.Println("Created directory:", path)
  return nil
}

/*
  Creates a file at target path, based on the contents of the provided io.Reader
*/
func createFile(fileName string, fileMode os.FileMode, reader io.Reader) error {
  parentDir := filepath.Dir(fileName)
  err := createDir(parentDir)
  if err != nil {
    return err
  }

  file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, fileMode)
  if err != nil {
    return err
  }
  if _, err := io.Copy(file, reader); err != nil {
    return err
  }
  file.Close()

  fmt.Println("Created file:", fileName)
  return nil
}

/*
  Create a symlink to sourceFile at symlinkFile
*/
func createSymlink(sourceFile string, symlinkFile string) {
  os.Symlink(sourceFile, symlinkFile)
}

/*
  Removes a directory at target path
*/
func removeDir(path string) error {
  return os.RemoveAll(path)
}

func removeFile(path string) error {
  return os.Remove(path)
}

/*
  Compare version a with version b,
  returns true if version b is equal or later than version a
*/
func compareVersions(a string, b string) (bool, error) {
  majorA, minorA, buildA, err := parseVersionString(a)
  if err != nil {
    return false, err
  }
  majorB, minorB, buildB, err := parseVersionString(b)
  if err != nil {
    return false, err
  }
  return majorB >= majorA && minorB >= minorA && buildB >= buildA, nil
}

/*
  Version format is major.minor.build
  Returns the integers in the version string
*/
func parseVersionString(version string) (int, int, int, error) {
  versionArr := strings.Split(version, ".")

  major, err := strconv.Atoi(versionArr[0])
  if err != nil {
    return 0, 0, 0, err
  }
  minor, err := strconv.Atoi(versionArr[1])
  if err != nil {
    return 0, 0, 0, err
  }
  build, err := strconv.Atoi(versionArr[2])
  if err != nil {
    return 0, 0, 0, err
  }

  return major, minor, build, nil
}

/*
  Generically handle error, print error, stacktrace and exit
*/
func handleErr(err error) {
  fmt.Println(err)
  debug.PrintStack()
  os.Exit(1)
}
