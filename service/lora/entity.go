package lora

type BaseResponse[T any] struct {
	Success bool   `json:"success"`
	Err     string `json:"err"`
	Code    int    `json:"code"`
	Data    T      `json:"data"`
}

type TranResponse struct {
	Config TrainConfigValues `json:"param"`
	Id     string            `json:"id"`
}

type TrainTask struct {
	ID          string                 `json:"id"`
	TaskType    string                 `json:"task_type"`
	Result      interface{}            `json:"result"`
	Status      string                 `json:"status"`
	Error       string                 `json:"error"`
	Output      []string               `json:"output"`
	Epoch       int                    `json:"epoch"`
	Steps       int                    `json:"steps"`
	TotalSteps  int                    `json:"total_steps"`
	TotalEpochs int                    `json:"total_epochs"`
	Config      map[string]interface{} `json:"config"`
}
type InfoResponse struct {
	Success bool   `json:"success"`
	Name    string `json:"name"`
}
type TrainConfigValues struct {
	Name                        string  `json:"name,omitempty"`
	PretrainedModelNameOrPath   string  `json:"pretrained_model_name_or_path"`
	V2                          bool    `json:"v2"`
	VParameterization           bool    `json:"v_parameterization"`
	TrainDataDir                string  `json:"train_data_dir"`
	OutputDir                   string  `json:"output_dir"`
	Resolution                  string  `json:"resolution"`
	LearningRate                float64 `json:"learning_rate"`
	LrScheduler                 string  `json:"lr_scheduler"`
	LrWarmupSteps               int     `json:"lr_warmup_steps,omitempty"`
	TrainBatchSize              int     `json:"train_batch_size"`
	Epoch                       string  `json:"epoch"`
	SaveEveryNEpochs            int     `json:"save_every_n_epochs,omitempty"`
	MixedPrecision              string  `json:"mixed_precision"`
	SavePrecision               string  `json:"save_precision"`
	Seed                        int     `json:"seed"`
	NumCpuThreadsPerProcess     int     `json:"num_cpu_threads_per_process"`
	CacheLatents                bool    `json:"cache_latents"`
	CaptionExtension            string  `json:"caption_extension"`
	EnableBucket                bool    `json:"enable_bucket"`
	GradientCheckpointing       bool    `json:"gradient_checkpointing"`
	FullFp16                    bool    `json:"full_fp16"`
	NoTokenPadding              bool    `json:"no_token_padding"`
	StopTextEncoderTrainingPct  float64 `json:"stop_text_encoder_training_pct"`
	Xformers                    bool    `json:"xformers"`
	SaveModelAs                 string  `json:"save_model_as"`
	ShuffleCaption              bool    `json:"shuffle_caption"`
	Resume                      string  `json:"resume,omitempty"`
	PriorLossWeight             float64 `json:"prior_loss_weight"`
	TextEncoderLr               float64 `json:"text_encoder_lr"`
	NetworkDim                  int     `json:"network_dim"`
	NetworkWeights              string  `json:"network_weights,omitempty"`
	ColorAug                    bool    `json:"color_aug"`
	FlipAug                     bool    `json:"flip_aug"`
	ClipSkip                    int     `json:"clip_skip"`
	GradientAccumulationSteps   int     `json:"gradient_accumulation_steps"`
	MemEffAttn                  bool    `json:"mem_eff_attn"`
	OutputName                  string  `json:"output_name"`
	ModelList                   string  `json:"model_list"`
	MaxTokenLength              int     `json:"max_token_length"`
	MaxTrainEpochs              int     `json:"max_train_epochs,omitempty"`
	MaxDataLoaderNWorkers       int     `json:"max_data_loader_n_workers"`
	NetworkAlpha                float64 `json:"network_alpha"`
	TrainingComment             string  `json:"training_comment"`
	KeepTokens                  int     `json:"keep_tokens,omitempty"`
	LrSchedulerNumCycles        int     `json:"lr_scheduler_num_cycles"`
	LrSchedulerPower            float64 `json:"lr_scheduler_power"`
	PersistentDataLoaderWorkers bool    `json:"persistent_data_loader_workers"`
	BucketNoUpscale             bool    `json:"bucket_no_upscale"`
	RandomCrop                  bool    `json:"random_crop"`
	BucketResoSteps             int     `json:"bucket_reso_steps"`
	CaptionDropoutEveryNEpochs  int     `json:"caption_dropout_every_n_epochs,omitempty"`
	CaptionTagDropoutRate       float64 `json:"caption_tag_dropout_rate"`
	Optimizer                   string  `json:"optimizer"`
	OptimizerArgs               string  `json:"optimizer_args,omitempty"`
	NoiseOffset                 string  `json:"noise_offset"`
	LoRAType                    string  `json:"LoRA_type"`
	ConvDim                     int     `json:"conv_dim"`
	ConvAlpha                   float64 `json:"conv_alpha"`
	SampleEveryNSteps           int     `json:"sample_every_n_steps,omitempty"`
	SampleEveryNEpochs          int     `json:"sample_every_n_epochs,omitempty"`
	SampleSampler               string  `json:"sample_sampler"`
	SamplePrompts               string  `json:"sample_prompts,omitempty"`
	AdditionalParameters        string  `json:"additional_parameters"`
	VaeBatchSize                int     `json:"vae_batch_size"`
	NetworkModule               string  `json:"network_module"`
	UNetLR                      float64 `json:"unet_lr"`
}
