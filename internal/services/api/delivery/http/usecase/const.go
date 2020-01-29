package handlers

const (
	ErrAuth           = "Required authorization"
	ErrAvatarNotFound = "User's avatar not found"
	ErrUserNotFound   = "User not found"

	// image handler errors
	ErrFailedImageSaveInService  = "Failed to save image in photo service"
	ErrFailedImageSaveInDatabase = "Failed to save image in our database"
	ErrInvalidFile               = "Invalid file provided"
	ErrInvalidFileSize           = "Invalid file size:%v. Maximum allowed size %v"
	ErrInvalidFileFormat         = "Invalid file format. You can upload an image in one of the following formats:"

	// image handler const
	FormFileName      = "file"
	ContentTypeHeader = "Content-Type"

	// session handler errors
	WrnFailedTokenCreate = "Failed create token in auth service"

	// user handler errors
	ErrUserAlreadyExist  = "User already exist"
	WrnFailedTokenDelete = "Failed delete token in auth service"

	// users handler errors
	ErrFailedPageCountGet = "Failed to get page count of users"
)
