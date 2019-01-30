package serviceapitests

import (
	"flag"
	"log"
	"os"
	"testing"
)

var TestBaseURL string

func TestMain(m *testing.M) {
	// parse and print command line flags
	flag.Parse()
	log.Printf("TestEnv=%s", TestConfig.TestEnv)
	log.Printf("Protocal=%s", TestConfig.Protocal)
	log.Printf("BasePath=%s", TestConfig.BasePath)
	log.Printf("Port=%s", TestConfig.Port)

	TestBaseURL = TestConfig.Protocal + "://" + TestConfig.BasePath + ":" + TestConfig.Port
	log.Printf("TestBaseUrl=%s", TestBaseURL)
	os.Exit(m.Run())
}
