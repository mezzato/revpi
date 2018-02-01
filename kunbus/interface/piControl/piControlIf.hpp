/*!
 *
 * Project: Pi Control
 * (C)    : KUNBUS GmbH, Wachhausstr. 5a, D-76227 Karlsruhe
 *
 *
 * \file piControlIf.c
 *
 * \brief PI Control Interface
 *
 *
 */

#ifndef PICONTROLIF_HPP_
#define PICONTROLIF_HPP_

#if !defined(__cplusplus)
#error must be compiled as CPP
#endif 

/******************************************************************************/
/********************************  Includes  **********************************/
/******************************************************************************/

#include <stdint.h>
#include <piControl.h>

/******************************************************************************/
/*********************************  Types  ************************************/
/******************************************************************************/

class piControl
{
private:
    int Handle_m;

    void Open(void);
    void Close(void);
    
public:
    piControl();
    ~piControl();
    
    int Reset(void);
    int Read(uint32_t Offset, uint32_t Length, uint8_t *pData);
    int Write(uint32_t Offset, uint32_t Length, uint8_t *pData);
    int GetDeviceInfo(SDeviceInfo *pDev);
    int GetDeviceInfoList(SDeviceInfo *pDev);
    int GetBitValue(SPIValue *pSpiValue);
    int SetBitValue(SPIValue *pSpiValue);
    int GetVariableInfo(SPIVariable *pSpiVariable);
    int FindVariable(const char *name);


};

#endif /* PICONTROLIF_HPP_ */
