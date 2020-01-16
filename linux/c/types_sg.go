/*
 * Author: fasion
 * Created time: 2019-08-28 18:08:39
 * Last Modified by: fasion
 * Last Modified time: 2019-08-29 14:37:00
 */

package c

/*
#include <scsi/sg.h>
*/
import "C"

const (
	SG_IO = C.SG_IO

	SG_DXFER_NONE = C.SG_DXFER_NONE
	SG_DXFER_TO_DEV = C.SG_DXFER_TO_DEV
	SG_DXFER_FROM_DEV = C.SG_DXFER_FROM_DEV
	SG_DXFER_TO_FROM_DEV = C.SG_DXFER_TO_FROM_DEV

	SG_GET_VERSION_NUM = C.SG_GET_VERSION_NUM
)


type SgIoHdrT_c = C.sg_io_hdr_t
