package reportsummary

import (
	"reflect"
	"testing"

	"github.com/kubescape/opa-utils/reporthandling/apis"
	"github.com/stretchr/testify/assert"
)

func TestSetStatus(t *testing.T) {
	status := []apis.ScanningStatus{}

	f := mockFrameworkSummaryFailExclude()
	status = append(status, f.GetStatus().Status())

	f.Status = apis.StatusUnknown
	for k, v := range f.Controls {
		status = append(status, v.GetStatus().Status())
		f.Controls[k] = v
	}

	i := 0
	assert.Equal(t, status[i], f.GetStatus().Status())
	for _, v := range f.Controls {
		i++
		assert.Equal(t, status[i], v.GetStatus().Status())
	}

}

func TestStatusInfoNotPresent(t *testing.T) {

	f := mockSummaryDetailsNoInnerStatus() // Status: skipped , InnerStatus empty
	for _, v := range f.Controls {
		status := v.GetStatus()
		assert.Equal(t, reflect.TypeOf(status), reflect.TypeOf(&apis.StatusInfo{}))
		assert.Equal(t, status.Status(), apis.InfoStatusSkipped)
		assert.Equal(t, status.Info(), "")
	}

}

func TestStatusEmpty(t *testing.T) {

	f := mockSummaryDetailsStatusEmpty()
	for _, v := range f.Controls {
		v.Status = apis.StatusIrrelevant
		status := v.GetStatus()
		assert.Equal(t, reflect.TypeOf(status), reflect.TypeOf(&apis.StatusInfo{}))
		assert.Equal(t, status.Status(), apis.StatusIrrelevant)
		assert.Equal(t, status.Info(), "")
	}

}

func TestStatusInfoSkipped(t *testing.T) {
	var status apis.ScanningStatus
	var info string

	f := mockSummaryDetailsStatusSkipped() // control -> status: "skipped", info: "no host sensor flag"

	for _, v := range f.Controls {
		status = v.GetStatus().Status()
		info = v.GetStatus().Info()
		assert.Equal(t, status, apis.InfoStatusSkipped)
		assert.Equal(t, info, "no host sensor flag")
	}

}

func TestStatusInfoIrelevant(t *testing.T) {
	var status apis.ScanningStatus
	var info string

	f := mockSummaryDetailsStatusIrrelevant() // control -> status: "irrelevant", info: "no k8s dashboard in cluster"

	for _, v := range f.Controls {
		status = v.GetStatus().Status()
		info = v.GetStatus().Info()
		assert.Equal(t, status, apis.StatusIrrelevant)
		assert.Equal(t, info, "no k8s dashboard in cluster")
	}

}

func TestFrameworkControlsSummariesCounters(t *testing.T) {
	f := mockFrameworkSummaryFailPass()
	f.ListControlsIDs().Skipped()
	assert.Equal(t, len(f.Controls), f.GetControls().NumberOfControls().All(), "invalid total control count")
	assert.Equal(t, len(f.GetControls().ListControlsIDs().Failed()), f.GetControls().NumberOfControls().Failed(), "invalid total failed control count")
	assert.Equal(t, len(f.GetControls().ListControlsIDs().Passed()), f.GetControls().NumberOfControls().Passed(), "invalid total passed control count")
	assert.Equal(t, len(f.GetControls().ListControlsIDs().Excluded()), f.GetControls().NumberOfControls().Excluded(), "invalid total excluded/warning control count")
	assert.Equal(t, len(f.GetControls().ListControlsIDs().Skipped()), f.GetControls().NumberOfControls().Skipped(), "invalid total skipped control count")
}

func TestFrameworkGettingSpecificControl(t *testing.T) {
	f := mockFrameworkSummaryFailPass()
	a := f.GetControls().GetControl(EControlCriteriaID, "1234")
	assert.Nil(t, a, "control id '1234' shouldn't exist")
	assert.Equal(t, "control-fail-pass", f.GetControls().GetControl(EControlCriteriaID, "C-0001").GetName(), "wrong control retrieved")
}
