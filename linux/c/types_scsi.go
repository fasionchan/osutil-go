// +build ignore

/*
 * Author: fasion
 * Created time: 2019-08-29 14:35:10
 * Last Modified by: fasion
 * Last Modified time: 2020-01-17 11:11:51
 */

package c

/*
#include <scsi/scsi.h>
*/
import "C"

const (
	TYPE_DISK           = C.TYPE_DISK
	TYPE_TAPE           = C.TYPE_TAPE
	TYPE_PROCESSOR      = C.TYPE_PROCESSOR
	TYPE_WORM           = C.TYPE_WORM
	TYPE_ROM            = C.TYPE_ROM
	TYPE_SCANNER        = C.TYPE_SCANNER
	TYPE_MOD            = C.TYPE_MOD
	TYPE_MEDIUM_CHANGER = C.TYPE_MEDIUM_CHANGER

	TYPE_NO_LUN = C.TYPE_NO_LUN
)
