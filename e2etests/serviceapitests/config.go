package serviceapitests

import "flag"

type testConfig struct {
	TestEnv  string
	Protocal string
	BasePath string
	Port     string
}

var (
	//TestConfig ...All Config var related to testing environments
	TestConfig testConfig
)

func init() {

	flag.StringVar(&TestConfig.TestEnv, "testEnv", "local", "Environment under test - such as qal, dev, e2e or preprod/prod")
	flag.StringVar(&TestConfig.Protocal, "protocal", "http", "Mode the Rest Service runs in - Secure/Insecure")
	flag.StringVar(&TestConfig.BasePath, "basePath", "localhost", "Base path of the API service under test")
	flag.StringVar(&TestConfig.Port, "port", "8011", "Port number of the API service under test")
}
