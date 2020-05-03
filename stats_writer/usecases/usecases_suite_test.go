package usecases_test

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/onsi/ginkgo/reporters"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var (
	testReporter gomock.TestReporter
	// MockController is the global mock controller to use in tests
	MockController = gomock.NewController(testReporter)
)

// RunTestSuite is used to launch the tests in a given suite
func RunTestSuite(name string, t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, name, []ginkgo.Reporter{junitReporter})
}

func TestMainStub(t *testing.T) {
	RunTestSuite("stats_writer usecases tests", t)
}
