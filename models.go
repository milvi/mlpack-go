package mlpack

/*
#include <capi/approx_kfn.h>
#include <capi/cf.h>
#include <capi/det.h>
#include <capi/fastmks.h>
#include <capi/gmm_train.h>
#include <capi/gmm_generate.h>
#include <capi/gmm_probability.h>
#include <capi/hmm_train.h>
#include <capi/hmm_generate.h>
#include <capi/hmm_loglik.h>
#include <capi/hmm_viterbi.h>
#include <capi/hoeffding_tree.h>
#include <capi/kde.h>
#include <capi/bayesian_linear_regression.h>
#include <capi/lars.h>
#include <capi/linear_svm.h>
#include <capi/local_coordinate_coding.h>
#include <capi/logistic_regression.h>
#include <capi/lsh.h>
#include <capi/nbc.h>
#include <capi/knn.h>
#include <capi/kfn.h>
#include <capi/perceptron.h>
#include <capi/preprocess_scale.h>
#include <capi/krann.h>
#include <capi/softmax_regression.h>
#include <capi/sparse_coding.h>
#include <capi/random_forest.h>
#include <capi/decision_tree.h>
#include <capi/adaboost.h>
#include <capi/linear_regression.h>

void mlpackDeleteAdaBoostModelPtr(void* value);

void mlpackDeleteApproxKFNModelPtr(void* value);

void mlpackDeleteBayesianLinearRegressionPtr(void* value);

void mlpackDeleteCFModelPtr(void* value);

void mlpackDeleteDecisionTreeModelPtr(void* value);

void mlpackDeleteDTreePtr(void* value);

void mlpackDeleteFastMKSModelPtr(void* value);

void mlpackDeleteGMMPtr(void* value);

void mlpackDeleteHMMModelPtr(void* value);

void mlpackDeleteHoeffdingTreeModelPtr(void* value);

void mlpackDeleteKDEModelPtr(void* value);

void mlpackDeleteLARSPtr(void* value);

void mlpackDeleteLinearRegressionPtr(void* value);

void mlpackDeleteLinearSVMModelPtr(void* value);

void mlpackDeleteLocalCoordinateCodingPtr(void* value);

void mlpackDeleteLogisticRegressionPtr(void* value);

void mlpackDeleteLSHSearchPtr(void* value);

void mlpackDeleteNBCModelPtr(void* value);

void mlpackDeleteKFNModelPtr(void* value);

void mlpackDeleteKNNModelPtr(void* value);

void mlpackDeletePerceptronModelPtr(void* value);

void mlpackDeleteScalingModelPtr(void* value);

void mlpackDeleteRandomForestModelPtr(void* value);

void mlpackDeleteRAModelPtr(void* value);

void mlpackDeleteSoftmaxRegressionPtr(void* value);

void mlpackDeleteSparseCodingPtr(void* value);
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type approxkfnModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *approxkfnModel) allocApproxKFNModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetApproxKFNModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteApproxKFNModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *approxkfnModel) getApproxKFNModel(params *params, identifier string) {
	m.allocApproxKFNModel(params, identifier)
}

func setApproxKFNModel(params *params,
	identifier string,
	ptr *approxkfnModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetApproxKFNModelPtr(params.mem,
		id, address)
}

type cfModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *cfModel) allocCFModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetCFModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteCFModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *cfModel) getCFModel(params *params, identifier string) {
	m.allocCFModel(params, identifier)
}

func setCFModel(params *params,
	identifier string,
	ptr *cfModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetCFModelPtr(params.mem,
		id, address)
}

type dTree struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *dTree) allocDTree(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetDTreePtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteDTreePtr(p) })
	runtime.KeepAlive(m)
}

func (m *dTree) getDTree(params *params, identifier string) {
	m.allocDTree(params, identifier)
}

func setDTree(params *params,
	identifier string,
	ptr *dTree) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetDTreePtr(params.mem,
		id, address)
}

type fastmksModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *fastmksModel) allocFastMKSModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetFastMKSModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteFastMKSModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *fastmksModel) getFastMKSModel(params *params, identifier string) {
	m.allocFastMKSModel(params, identifier)
}

func setFastMKSModel(params *params,
	identifier string,
	ptr *fastmksModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetFastMKSModelPtr(params.mem,
		id, address)
}

type gmm struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *gmm) allocGMM(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetGMMPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteGMMPtr(p) })
	runtime.KeepAlive(m)
}

func (m *gmm) getGMM(params *params, identifier string) {
	m.allocGMM(params, identifier)
}

func setGMM(params *params,
	identifier string,
	ptr *gmm) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetGMMPtr(params.mem,
		id, address)
}

type hmmModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *hmmModel) allocHMMModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetHMMModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteHMMModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *hmmModel) getHMMModel(params *params, identifier string) {
	m.allocHMMModel(params, identifier)
}

func setHMMModel(params *params,
	identifier string,
	ptr *hmmModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetHMMModelPtr(params.mem,
		id, address)
}

type hoeffdingTreeModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *hoeffdingTreeModel) allocHoeffdingTreeModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetHoeffdingTreeModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteHoeffdingTreeModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *hoeffdingTreeModel) getHoeffdingTreeModel(params *params, identifier string) {
	m.allocHoeffdingTreeModel(params, identifier)
}

func setHoeffdingTreeModel(params *params,
	identifier string,
	ptr *hoeffdingTreeModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetHoeffdingTreeModelPtr(params.mem,
		id, address)
}

type kdeModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *kdeModel) allocKDEModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetKDEModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteKDEModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *kdeModel) getKDEModel(params *params, identifier string) {
	m.allocKDEModel(params, identifier)
}

func setKDEModel(params *params,
	identifier string,
	ptr *kdeModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetKDEModelPtr(params.mem,
		id, address)
}

type bayesianLinearRegression struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *bayesianLinearRegression) allocBayesianLinearRegression(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetBayesianLinearRegressionPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteBayesianLinearRegressionPtr(p) })
	runtime.KeepAlive(m)
}

func (m *bayesianLinearRegression) getBayesianLinearRegression(params *params, identifier string) {
	m.allocBayesianLinearRegression(params, identifier)
}

func setBayesianLinearRegression(params *params,
	identifier string,
	ptr *bayesianLinearRegression) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetBayesianLinearRegressionPtr(params.mem,
		id, address)
}

type lars struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *lars) allocLARS(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetLARSPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteLARSPtr(p) })
	runtime.KeepAlive(m)
}

func (m *lars) getLARS(params *params, identifier string) {
	m.allocLARS(params, identifier)
}

func setLARS(params *params,
	identifier string,
	ptr *lars) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetLARSPtr(params.mem,
		id, address)
}

type linearsvmModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *linearsvmModel) allocLinearSVMModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetLinearSVMModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteLinearSVMModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *linearsvmModel) getLinearSVMModel(params *params, identifier string) {
	m.allocLinearSVMModel(params, identifier)
}

func setLinearSVMModel(params *params,
	identifier string,
	ptr *linearsvmModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetLinearSVMModelPtr(params.mem,
		id, address)
}

type localCoordinateCoding struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *localCoordinateCoding) allocLocalCoordinateCoding(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetLocalCoordinateCodingPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteLocalCoordinateCodingPtr(p) })
	runtime.KeepAlive(m)
}

func (m *localCoordinateCoding) getLocalCoordinateCoding(params *params, identifier string) {
	m.allocLocalCoordinateCoding(params, identifier)
}

func setLocalCoordinateCoding(params *params,
	identifier string,
	ptr *localCoordinateCoding) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetLocalCoordinateCodingPtr(params.mem,
		id, address)
}

type logisticRegression struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *logisticRegression) allocLogisticRegression(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetLogisticRegressionPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteLogisticRegressionPtr(p) })
	runtime.KeepAlive(m)
}

func (m *logisticRegression) getLogisticRegression(params *params, identifier string) {
	m.allocLogisticRegression(params, identifier)
}

func setLogisticRegression(params *params,
	identifier string,
	ptr *logisticRegression) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetLogisticRegressionPtr(params.mem,
		id, address)
}

type lshSearch struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *lshSearch) allocLSHSearch(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetLSHSearchPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteLSHSearchPtr(p) })
	runtime.KeepAlive(m)
}

func (m *lshSearch) getLSHSearch(params *params, identifier string) {
	m.allocLSHSearch(params, identifier)
}

func setLSHSearch(params *params,
	identifier string,
	ptr *lshSearch) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetLSHSearchPtr(params.mem,
		id, address)
}

type nbcModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *nbcModel) allocNBCModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetNBCModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteNBCModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *nbcModel) getNBCModel(params *params, identifier string) {
	m.allocNBCModel(params, identifier)
}

func setNBCModel(params *params,
	identifier string,
	ptr *nbcModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetNBCModelPtr(params.mem,
		id, address)
}

type knnModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *knnModel) allocKNNModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetKNNModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteKNNModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *knnModel) getKNNModel(params *params, identifier string) {
	m.allocKNNModel(params, identifier)
}

func setKNNModel(params *params,
	identifier string,
	ptr *knnModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetKNNModelPtr(params.mem,
		id, address)
}

type kfnModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *kfnModel) allocKFNModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetKFNModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteKFNModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *kfnModel) getKFNModel(params *params, identifier string) {
	m.allocKFNModel(params, identifier)
}

func setKFNModel(params *params,
	identifier string,
	ptr *kfnModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetKFNModelPtr(params.mem,
		id, address)
}

type perceptronModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *perceptronModel) allocPerceptronModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetPerceptronModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeletePerceptronModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *perceptronModel) getPerceptronModel(params *params, identifier string) {
	m.allocPerceptronModel(params, identifier)
}

func setPerceptronModel(params *params,
	identifier string,
	ptr *perceptronModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetPerceptronModelPtr(params.mem,
		id, address)
}

type scalingModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *scalingModel) allocScalingModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetScalingModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteScalingModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *scalingModel) getScalingModel(params *params, identifier string) {
	m.allocScalingModel(params, identifier)
}

func setScalingModel(params *params,
	identifier string,
	ptr *scalingModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetScalingModelPtr(params.mem,
		id, address)
}

type raModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *raModel) allocRAModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetRAModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteRAModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *raModel) getRAModel(params *params, identifier string) {
	m.allocRAModel(params, identifier)
}

func setRAModel(params *params,
	identifier string,
	ptr *raModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetRAModelPtr(params.mem,
		id, address)
}

type softmaxRegression struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *softmaxRegression) allocSoftmaxRegression(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetSoftmaxRegressionPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteSoftmaxRegressionPtr(p) })
	runtime.KeepAlive(m)
}

func (m *softmaxRegression) getSoftmaxRegression(params *params, identifier string) {
	m.allocSoftmaxRegression(params, identifier)
}

func setSoftmaxRegression(params *params,
	identifier string,
	ptr *softmaxRegression) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetSoftmaxRegressionPtr(params.mem,
		id, address)
}

type sparseCoding struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *sparseCoding) allocSparseCoding(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetSparseCodingPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteSparseCodingPtr(p) })
	runtime.KeepAlive(m)
}

func (m *sparseCoding) getSparseCoding(params *params, identifier string) {
	m.allocSparseCoding(params, identifier)
}

func setSparseCoding(params *params,
	identifier string,
	ptr *sparseCoding) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetSparseCodingPtr(params.mem,
		id, address)
}

type randomForestModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *randomForestModel) allocRandomForestModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetRandomForestModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteRandomForestModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *randomForestModel) getRandomForestModel(params *params, identifier string) {
	m.allocRandomForestModel(params, identifier)
}

func setRandomForestModel(params *params,
	identifier string,
	ptr *randomForestModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetRandomForestModelPtr(params.mem,
		id, address)
}

type decisionTreeModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *decisionTreeModel) allocDecisionTreeModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetDecisionTreeModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteDecisionTreeModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *decisionTreeModel) getDecisionTreeModel(params *params, identifier string) {
	m.allocDecisionTreeModel(params, identifier)
}

func setDecisionTreeModel(params *params,
	identifier string,
	ptr *decisionTreeModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetDecisionTreeModelPtr(params.mem,
		id, address)
}

type adaBoostModel struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *adaBoostModel) allocAdaBoostModel(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetAdaBoostModelPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteAdaBoostModelPtr(p) })
	runtime.KeepAlive(m)
}

func (m *adaBoostModel) getAdaBoostModel(params *params, identifier string) {
	m.allocAdaBoostModel(params, identifier)
}

func setAdaBoostModel(params *params,
	identifier string,
	ptr *adaBoostModel) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetAdaBoostModelPtr(params.mem,
		id, address)
}

type linearRegression struct {
	owner *modelOwner
	mem   unsafe.Pointer
}

func (m *linearRegression) allocLinearRegression(params *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	m.mem = C.mlpackGetLinearRegressionPtr(params.mem,
		id)
	checkNative()
	m.owner = adoptModel(params, m.mem, func(p unsafe.Pointer) { C.mlpackDeleteLinearRegressionPtr(p) })
	runtime.KeepAlive(m)
}

func (m *linearRegression) getLinearRegression(params *params, identifier string) {
	m.allocLinearRegression(params, identifier)
}

func setLinearRegression(params *params,
	identifier string,
	ptr *linearRegression) {
	address := ptr.owner.address()
	params.refs = append(params.refs, ptr.owner)
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetLinearRegressionPtr(params.mem,
		id, address)
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *adaBoostModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *approxkfnModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *bayesianLinearRegression) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *cfModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *decisionTreeModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *dTree) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *fastmksModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *gmm) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *hmmModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *hoeffdingTreeModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *kdeModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *lars) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *linearRegression) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *linearsvmModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *localCoordinateCoding) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *logisticRegression) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *lshSearch) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *nbcModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *kfnModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *knnModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *perceptronModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *scalingModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *randomForestModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *raModel) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *softmaxRegression) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}

// Close libère le modèle ; les alias partagent le même propriétaire.
func (m *sparseCoding) Close() error {
	if m == nil {
		return nil
	}
	return m.owner.Close()
}
