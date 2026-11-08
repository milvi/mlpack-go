/**
 * @file bindings/go/capi/arma_util.h
 * @author Yasmine Dumouchl
 * @author Yashwant Singh
 *
 * Header file for cgo to call C functions from go.
 *
 * mlpack is free software; you may redistribute it and/or modify it under the
 * terms of the 3-clause BSD license.  You should have received a copy of the
 * 3-clause BSD license along with mlpack.  If not, see
 * http://www.opensource.org/licenses/BSD-3-Clause for more information.
 */
#ifndef MLPACK_BINDINGS_GO_MLPACK_ARMAUTIL_H
#define MLPACK_BINDINGS_GO_MLPACK_ARMAUTIL_H

#include <stdint.h>
#include <stddef.h>
#include <stdbool.h>

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

// Copie sans détachement ; kind : mat, umat, row, urow, col, ucol, tuple (0 à 6).
bool mlpackCopyParamMat(void* params, const char* identifier,
                       double* destination, const size_t length, const int kind);


/**
 * Pass Gonum Dense pointer and wrap an Armadillo mat around it.
 */
void mlpackToArmaMat(void* params,
                     const char* identifier,
                     double* mat,
                     const size_t row,
                     const size_t col,
                     bool transpose);

/**
 * Pass Gonum Dense pointer and wrap an Armadillo mat around it.
 */
void mlpackToArmaUmat(void* params,
                      const char* identifier,
                      double* mat,
                      const size_t row,
                      const size_t col);

/**
 * Pass Gonum VecDense pointer and wrap an Armadillo rowvec around it.
 */
void mlpackToArmaRow(void* params,
                     const char* identifier,
                     double* rowvec,
                     const size_t elem);

/**
 * Pass Gonum VecDense pointer and wrap an Armadillo rowvec around it.
 */
void mlpackToArmaUrow(void* params,
                      const char* identifier,
                      double* rowvec,
                      const size_t elem);

/**
 * Pass Gonum VecDense pointer and wrap an Armadillo colvec around it.
 */
void mlpackToArmaCol(void* params,
                     const char* identifier,
                     double* colvec,
                     const size_t elem);

/**
 * Pass Gonum VecDense pointer and wrap an Armadillo colvec around it.
 */
void mlpackToArmaUcol(void* params,
                      const char* identifier,
                      double* colvec,
                      const size_t elem);

/**
 * Return the memory pointer of an Armadillo mat object.
 */

/**
 * Return the memory pointer of an Armadillo umat object.
 */

/**
 * Return the memory pointer of an Armadillo row object.
 */

/**
 * Return the memory pointer of an Armadillo urow object.
 */

/**
 * Return the memory pointer of an Armadillo col object.
 */

/**
 * Return the memory pointer of an Armadillo ucol object.
 */

/**
 * Return the number of rows in a Armadillo mat.
 */
size_t mlpackNumRowMat(void* params, const char* identifier);

/**
 * Return the number of columns in an Armadillo mat.
 */
size_t mlpackNumColMat(void* params, const char* identifier);

/**
 * Return the number of elements in an Armadillo mat.
 */
size_t mlpackNumElemMat(void* params, const char* identifier);

/**
 * Return the number of rows in an Armadillo umat.
 */
size_t mlpackNumRowUmat(void* params, const char* identifier);

/**
 * Return the number of columns in an Armadillo umat.
 */
size_t mlpackNumColUmat(void* params, const char* identifier);

/**
 * Return the number of elements in an Armadillo umat.
 */
size_t mlpackNumElemUmat(void* params, const char* identifier);

/**
 * Return the number of elements in an Armadillo row.
 */
size_t mlpackNumElemRow(void* params, const char* identifier);

/**
 * Return the number of elements in an Armadillo urow.
 */
size_t mlpackNumElemUrow(void* params, const char* identifier);

/**
 * Return the number of elements in an Armadillo col.
 */
size_t mlpackNumElemCol(void* params, const char* identifier);

/**
 * Return the number of elements in an Armadillo ucol.
 */
size_t mlpackNumElemUcol(void* params, const char* identifier);

/**
 * Call IO::SetParam<std::tuple<DatasetInfo, arma::mat>>().
 */
void mlpackToArmaMatWithInfo(void* params,
                             const char* identifier,
                             const bool* dimensions,
                             double* memptr,
                             const size_t rows,
                             const size_t cols);

/**
 * Get the number of elements in a matrix with DatasetInfo parameter.
 */
size_t mlpackArmaMatWithInfoElements(void* params, const char* identifier);

/**
 * Get the number of rows in a matrix with DatasetInfo parameter.
 */
size_t mlpackArmaMatWithInfoRows(void* params, const char* identifier);

/**
 * Get the number of columns in a matrix with DatasetInfo parameter.
 */
size_t mlpackArmaMatWithInfoCols(void* params, const char* identifier);

/**
 * Get a pointer to the memory of the matrix.  The calling function is expected
 * to own the memory.
 */


#if defined(__cplusplus) || defined(c_plusplus)
}
#endif

#endif
