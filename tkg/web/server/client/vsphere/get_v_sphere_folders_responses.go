// Code generated by go-swagger; DO NOT EDIT.

package vsphere

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/vmware-tanzu/tanzu-framework/tkg/web/server/models"
)

// GetVSphereFoldersReader is a Reader for the GetVSphereFolders structure.
type GetVSphereFoldersReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetVSphereFoldersReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetVSphereFoldersOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetVSphereFoldersBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewGetVSphereFoldersUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetVSphereFoldersInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetVSphereFoldersOK creates a GetVSphereFoldersOK with default headers values
func NewGetVSphereFoldersOK() *GetVSphereFoldersOK {
	return &GetVSphereFoldersOK{}
}

/*GetVSphereFoldersOK handles this case with default header values.

Successful retrieval of vSphere folders
*/
type GetVSphereFoldersOK struct {
	Payload []*models.VSphereFolder
}

func (o *GetVSphereFoldersOK) Error() string {
	return fmt.Sprintf("[GET /api/providers/vsphere/folders][%d] getVSphereFoldersOK  %+v", 200, o.Payload)
}

func (o *GetVSphereFoldersOK) GetPayload() []*models.VSphereFolder {
	return o.Payload
}

func (o *GetVSphereFoldersOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetVSphereFoldersBadRequest creates a GetVSphereFoldersBadRequest with default headers values
func NewGetVSphereFoldersBadRequest() *GetVSphereFoldersBadRequest {
	return &GetVSphereFoldersBadRequest{}
}

/*GetVSphereFoldersBadRequest handles this case with default header values.

Bad request
*/
type GetVSphereFoldersBadRequest struct {
	Payload *models.Error
}

func (o *GetVSphereFoldersBadRequest) Error() string {
	return fmt.Sprintf("[GET /api/providers/vsphere/folders][%d] getVSphereFoldersBadRequest  %+v", 400, o.Payload)
}

func (o *GetVSphereFoldersBadRequest) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetVSphereFoldersBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetVSphereFoldersUnauthorized creates a GetVSphereFoldersUnauthorized with default headers values
func NewGetVSphereFoldersUnauthorized() *GetVSphereFoldersUnauthorized {
	return &GetVSphereFoldersUnauthorized{}
}

/*GetVSphereFoldersUnauthorized handles this case with default header values.

Incorrect credentials
*/
type GetVSphereFoldersUnauthorized struct {
	Payload *models.Error
}

func (o *GetVSphereFoldersUnauthorized) Error() string {
	return fmt.Sprintf("[GET /api/providers/vsphere/folders][%d] getVSphereFoldersUnauthorized  %+v", 401, o.Payload)
}

func (o *GetVSphereFoldersUnauthorized) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetVSphereFoldersUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetVSphereFoldersInternalServerError creates a GetVSphereFoldersInternalServerError with default headers values
func NewGetVSphereFoldersInternalServerError() *GetVSphereFoldersInternalServerError {
	return &GetVSphereFoldersInternalServerError{}
}

/*GetVSphereFoldersInternalServerError handles this case with default header values.

Internal server error
*/
type GetVSphereFoldersInternalServerError struct {
	Payload *models.Error
}

func (o *GetVSphereFoldersInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/providers/vsphere/folders][%d] getVSphereFoldersInternalServerError  %+v", 500, o.Payload)
}

func (o *GetVSphereFoldersInternalServerError) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetVSphereFoldersInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}