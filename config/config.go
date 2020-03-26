package config

// Configuration for app
type Configuration struct {
	CameraStillPicURL    string `mapstructure:"CAMERA_STILL_PIC_URL"`
	CameraMotionURL      string `mapstructure:"CAMERA_MOTION_URL"`
	CameraLoginURL       string `mapstructure:"CAMERA_LOGIN_URL"`
	CameraUsername       string `mapstructure:"CAMERA_USERNAME"`
	CameraPassword       string `mapstructure:"CAMERA_PASSWORD"`
	CaptureFolder        string `mapstructure:"CAPTURE_FOLDER"`
	HTTPRetryCount       int    `mapstructure:"HTTP_RETRY_COUNT"`
	PhotosInSet          int    `mapstructure:"NUMBER_PHOTOS_IN_SET"`
	TimeoutValue         int    `mapstructure:"TIMEOUT"`
	ProjectID            string `mapstructure:"CUSTOM_VISION_PROJECT_ID"`
	ProjectIDDirection   string `mapstructure:"CUSTOM_VISION_PROJECT_DIRECTION_ID"`
	PredictionKey        string `mapstructure:"CUSTOM_VISION_PREDICTION_KEY"`
	EndpointURL          string `mapstructure:"CUSTOM_VISION_ENDPOINT"`
	IterationID          string `mapstructure:"CUSTOM_VISION_ITERATION_ID"`
	IterationIDDirection string `mapstructure:"CUSTOM_VISION_ITERATION_DIRECTION_ID"`
	ProcessedFolder      string `mapstructure:"PROCESSED_FOLDER"`
	FirebaseCredentials  string `mapstructure:"GOOGLE_FIREBASE_CREDENTIAL_FILE"`
	FirestoreCollection  string `mapstructure:"GOOGLE_FIRESTORE_COLLECTION"`
	StorageBucket        string `mapstructure:"STORAGE_BUCKET"`
	StorageFolder        string `mapstructure:"STORAGE_FOLDER"`
}
