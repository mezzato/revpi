/*!
 *
 * Project: Pi Control
 * Demo source code for usage of piControl driver
 *
 *   Copyright (C) 2016 : KUNBUS GmbH, Heerweg 15C, 73370 Denkendorf, Germany
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <http://www.gnu.org/licenses/>. 
 *
 *
 * \file piControlIf.c
 *
 * \brief PI Control Interface
 *
 *
 */

/******************************************************************************/
/********************************  Includes  **********************************/
/******************************************************************************/

#include "piControlIf.hpp"

#include "piControl.h"

#include <sys/types.h>
#include <sys/stat.h>
#include <sys/ioctl.h>
#include <fcntl.h>
#include <unistd.h>
#include <errno.h>

#include <stdio.h>
#include <string.h>


piControl::piControl(void)
{
    Handle_m = -1;
}


piControl::~piControl(void)
{
    Close();
}


/***********************************************************************************/
/*!
 * @brief Open Pi Control Interface
 *
 * Initialize the Pi Control Interface
 *
 ************************************************************************************/
void piControl::Open(void)
{
    /* open handle if needed */
    if (Handle_m < 0)
    {
        Handle_m = open(PICONTROL_DEVICE, O_RDWR);
    }
}


/***********************************************************************************/
/*!
 * @brief Close Pi Control Interface
 *
 * Clsoe the Pi Control Interface
 *
 ************************************************************************************/
void piControl::Close(void)
{
    /* open handle if needed */
    if (Handle_m > 0)
    {
        close(Handle_m);
        Handle_m = -1;
    }
}


/***********************************************************************************/
/*!
 * @brief Reset Pi Control Interface
 *
 * Initialize the Pi Control Interface
 *
 ************************************************************************************/
int piControl::Reset(void)
{
    Open();

    if (Handle_m < 0)
        return -18;

    // do some ioctls
    if (ioctl(Handle_m, KB_RESET, NULL) < 0)
        return errno;
    
    return 0;
}


/***********************************************************************************/
/*!
 * @brief Get Processdata
 *
 * Gets Processdata from a specific position
 *
 * @param[in]   Offset
 * @param[in]   Length
 * @param[out]  pData
 *
 * @return Number of Bytes read or error if negative
 *
 ************************************************************************************/
int piControl::Read(uint32_t Offset, uint32_t Length, uint8_t *pData)
{
    int BytesRead = 0;

    Open();

    if (Handle_m < 0)
        return -1;

    /* seek */
    if (lseek(Handle_m, Offset, SEEK_SET) < 0)
    {
        return -2;
    }

    /* read */
    BytesRead = read(Handle_m, pData, Length);
    if (BytesRead < 0)
    {
        return -3;
    }

    return BytesRead;
}


/***********************************************************************************/
/*!
 * @brief Set Processdata
 *
 * Writes Processdata at a specific position
 *
 * @param[in]   Offset
 * @param[in]   Length
 * @param[out]  pData
 *
 * @return Number of Bytes written or error if negative
 *
 ************************************************************************************/
int piControl::Write(uint32_t Offset, uint32_t Length, uint8_t *pData)
{
    int BytesWritten = 0;

    Open();

    if (Handle_m < 0)
        return -1;

    /* seek */
    if (lseek(Handle_m, Offset, SEEK_SET) < 0)
    {
        return -2;
    }

    /* Write */
    BytesWritten = write(Handle_m, pData, Length);
    if (BytesWritten < 0)
    {
        return -3;
    }

    return BytesWritten;
}


/***********************************************************************************/
/*!
 * @brief Get Device Info
 *
 * Get Description of connected devices.
 *
 * @param[in/out]   Pointer to an array of 20 entries of type SDeviceInfo.
 *
 * @return Number of detected devices
 *
 ************************************************************************************/
int piControl::GetDeviceInfo(SDeviceInfo *pDev)
{
    Open();

    if (Handle_m < 0)
        return -1;

    return ioctl(Handle_m, KB_GET_DEVICE_INFO, pDev);
}

int piControl::GetDeviceInfoList(SDeviceInfo *pDev)
{
    Open();

    if (Handle_m < 0)
        return -1;

    return ioctl(Handle_m, KB_GET_DEVICE_INFO_LIST, pDev);
}

/***********************************************************************************/
/*!
 * @brief Get Bit Value
 *
 * Get the value of one bit in the process image.
 *
 * @param[in/out]   Pointer to SPIValue.
 *
 * @return 0 or error if negative
 *
 ************************************************************************************/
int piControl::GetBitValue(SPIValue *pSpiValue)
{
    Open();

    if (Handle_m < 0)
        return -1;

    return ioctl(Handle_m, KB_GET_VALUE, pSpiValue);
}

/***********************************************************************************/
/*!
 * @brief Set Bit Value
 *
 * Set the value of one bit in the process image.
 *
 * @param[in/out]   Pointer to SPIValue.
 *
 * @return 0 or error if negative
 *
 ************************************************************************************/
int piControl::SetBitValue(SPIValue *pSpiValue)
{
    Open();

    if (Handle_m < 0)
        return -1;

    return ioctl(Handle_m, KB_SET_VALUE, pSpiValue);
}

/***********************************************************************************/
/*!
 * @brief Get Variable Info
 *
 * Get the info for a variable.
 *
 * @param[in/out]   Pointer to SPIVariable.
 *
 * @return 0 or error if negative
 *
 ************************************************************************************/
int piControl::GetVariableInfo(SPIVariable *pSpiVariable)
{
    Open();

    if (Handle_m < 0)
        return -1;

    return ioctl(Handle_m, KB_FIND_VARIABLE, pSpiVariable);
}

/***********************************************************************************/
/*!
 * @brief Get Variable offset by name
 *
 * Get the offset of a variable in the process image. This does NOT work for variable of type bool.
 *
 * @param[in]   pointer to string with name of variable
 *
 * @return      >= 0    offset
                < 0     in case of error
 *
 ************************************************************************************/
int piControl::FindVariable(const char *name)
{
    int ret;
    SPIVariable var;
        
    Open();

    strncpy(var.strVarName, name, sizeof(var.strVarName));
    var.strVarName[sizeof(var.strVarName) - 1] = 0;
        
    ret = ioctl(Handle_m, KB_FIND_VARIABLE, &var);
    if (ret < 0)
    {
        //printf("could not find variable '%s' in configuration.\n", var.strVarName);
    }
    else
    {
        //printf("Variable '%s' is at offset %d and %d bits long\n", var.strVarName, var.i16uAddress, var.i16uLength);
        ret = var.i16uAddress;
    }
    return ret;     
}
