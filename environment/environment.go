package environment

// PubKeyFile : is the public key file name
const PubKeyFile string = "rsa_public_key.pub"

// PrivKeyFile : is the private key file name
const PrivKeyFile string = "rsa_private_key.priv"

// EncryptedKeyFile : is the encrypted key file
const EncryptedKeyFile string = "safe_key"

// HomeFolder : is the HomeFolder name
const HomeFolder string = "You"

// UserPath : is the path of the user who executed the exe
var UserPath string

// EncryptionRootPath : is the path of the root folder of which all sub directories will be encrypted
var EncryptionRootPath string

// HomePath : is the path of the software home
var HomePath string

// ExePath : is the path of the software itself
var ExePath string
