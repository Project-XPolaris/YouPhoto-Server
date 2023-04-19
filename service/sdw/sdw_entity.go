package sdw

type Text2ImageResponse struct {
	Images     []string             `json:"images"`
	Parameters Text2ImageParameters `json:"parameters"`
	Info       string               `json:"info"`
}
type AlwaysonScripts struct {
}
type Text2ImageParameters struct {
	EnableHr                          bool            `json:"enable_hr"`
	DenoisingStrength                 float64         `json:"denoising_strength"`
	FirstphaseWidth                   int             `json:"firstphase_width"`
	FirstphaseHeight                  int             `json:"firstphase_height"`
	HrScale                           float64         `json:"hr_scale"`
	HrUpscaler                        string          `json:"hr_upscaler"`
	HrSecondPassSteps                 float64         `json:"hr_second_pass_steps"`
	HrResizeX                         int             `json:"hr_resize_x"`
	HrResizeY                         int             `json:"hr_resize_y"`
	Prompt                            string          `json:"prompt"`
	Styles                            interface{}     `json:"styles"`
	Seed                              int             `json:"seed"`
	Subseed                           int             `json:"subseed"`
	SubseedStrength                   int             `json:"subseed_strength"`
	SeedResizeFromH                   int             `json:"seed_resize_from_h"`
	SeedResizeFromW                   int             `json:"seed_resize_from_w"`
	SamplerName                       interface{}     `json:"sampler_name"`
	BatchSize                         int             `json:"batch_size"`
	NIter                             int             `json:"n_iter"`
	Steps                             int             `json:"steps"`
	CfgScale                          float64         `json:"cfg_scale"`
	Width                             int             `json:"width"`
	Height                            int             `json:"height"`
	RestoreFaces                      bool            `json:"restore_faces"`
	Tiling                            bool            `json:"tiling"`
	DoNotSaveSamples                  bool            `json:"do_not_save_samples"`
	DoNotSaveGrid                     bool            `json:"do_not_save_grid"`
	NegativePrompt                    string          `json:"negative_prompt"`
	Eta                               interface{}     `json:"eta"`
	SChurn                            float64         `json:"s_churn"`
	STmax                             interface{}     `json:"s_tmax"`
	STmin                             float64         `json:"s_tmin"`
	SNoise                            float64         `json:"s_noise"`
	OverrideSettings                  interface{}     `json:"override_settings"`
	OverrideSettingsRestoreAfterwards bool            `json:"override_settings_restore_afterwards"`
	ScriptArgs                        []interface{}   `json:"script_args"`
	SamplerIndex                      string          `json:"sampler_index"`
	ScriptName                        interface{}     `json:"script_name"`
	SendImages                        bool            `json:"send_images"`
	SaveImages                        bool            `json:"save_images"`
	AlwaysonScripts                   AlwaysonScripts `json:"alwayson_scripts"`
}

type Options struct {
	SamplesSave                        bool          `json:"samples_save"`
	SamplesFormat                      string        `json:"samples_format"`
	SamplesFilenamePattern             string        `json:"samples_filename_pattern"`
	SaveImagesAddNumber                bool          `json:"save_images_add_number"`
	GridSave                           bool          `json:"grid_save"`
	GridFormat                         string        `json:"grid_format"`
	GridExtendedFilename               bool          `json:"grid_extended_filename"`
	GridOnlyIfMultiple                 bool          `json:"grid_only_if_multiple"`
	GridPreventEmptySpots              bool          `json:"grid_prevent_empty_spots"`
	NRows                              float64       `json:"n_rows"`
	EnablePnginfo                      bool          `json:"enable_pnginfo"`
	SaveTxt                            bool          `json:"save_txt"`
	SaveImagesBeforeFaceRestoration    bool          `json:"save_images_before_face_restoration"`
	SaveImagesBeforeHighresFix         bool          `json:"save_images_before_highres_fix"`
	SaveImagesBeforeColorCorrection    bool          `json:"save_images_before_color_correction"`
	JpegQuality                        float64       `json:"jpeg_quality"`
	WebpLossless                       bool          `json:"webp_lossless"`
	ExportFor4Chan                     bool          `json:"export_for_4chan"`
	ImgDownscaleThreshold              float64       `json:"img_downscale_threshold"`
	TargetSideLength                   float64       `json:"target_side_length"`
	ImgMaxSizeMp                       float64       `json:"img_max_size_mp"`
	UseOriginalNameBatch               bool          `json:"use_original_name_batch"`
	UseUpscalerNameAsSuffix            bool          `json:"use_upscaler_name_as_suffix"`
	SaveSelectedOnly                   bool          `json:"save_selected_only"`
	DoNotAddWatermark                  bool          `json:"do_not_add_watermark"`
	TempDir                            string        `json:"temp_dir"`
	CleanTempDirAtStart                bool          `json:"clean_temp_dir_at_start"`
	OutdirSamples                      string        `json:"outdir_samples"`
	OutdirTxt2ImgSamples               string        `json:"outdir_txt2img_samples"`
	OutdirImg2ImgSamples               string        `json:"outdir_img2img_samples"`
	OutdirExtrasSamples                string        `json:"outdir_extras_samples"`
	OutdirGrids                        string        `json:"outdir_grids"`
	OutdirTxt2ImgGrids                 string        `json:"outdir_txt2img_grids"`
	OutdirImg2ImgGrids                 string        `json:"outdir_img2img_grids"`
	OutdirSave                         string        `json:"outdir_save"`
	SaveToDirs                         bool          `json:"save_to_dirs"`
	GridSaveToDirs                     bool          `json:"grid_save_to_dirs"`
	UseSaveToDirsForUI                 bool          `json:"use_save_to_dirs_for_ui"`
	DirectoriesFilenamePattern         string        `json:"directories_filename_pattern"`
	DirectoriesMaxPromptWords          float64       `json:"directories_max_prompt_words"`
	ESRGANTile                         float64       `json:"ESRGAN_tile"`
	ESRGANTileOverlap                  float64       `json:"ESRGAN_tile_overlap"`
	RealesrganEnabledModels            []string      `json:"realesrgan_enabled_models"`
	UpscalerForImg2Img                 interface{}   `json:"upscaler_for_img2img"`
	FaceRestorationModel               string        `json:"face_restoration_model"`
	CodeFormerWeight                   float64       `json:"code_former_weight"`
	FaceRestorationUnload              bool          `json:"face_restoration_unload"`
	ShowWarnings                       bool          `json:"show_warnings"`
	MemmonPollRate                     float64       `json:"memmon_poll_rate"`
	SamplesLogStdout                   bool          `json:"samples_log_stdout"`
	MultipleTqdm                       bool          `json:"multiple_tqdm"`
	PrintHypernetExtra                 bool          `json:"print_hypernet_extra"`
	UnloadModelsWhenTraining           bool          `json:"unload_models_when_training"`
	PinMemory                          bool          `json:"pin_memory"`
	SaveOptimizerState                 bool          `json:"save_optimizer_state"`
	SaveTrainingSettingsToTxt          bool          `json:"save_training_settings_to_txt"`
	DatasetFilenameWordRegex           string        `json:"dataset_filename_word_regex"`
	DatasetFilenameJoinString          string        `json:"dataset_filename_join_string"`
	TrainingImageRepeatsPerEpoch       float64       `json:"training_image_repeats_per_epoch"`
	TrainingWriteCsvEvery              float64       `json:"training_write_csv_every"`
	TrainingXattentionOptimizations    bool          `json:"training_xattention_optimizations"`
	TrainingEnableTensorboard          bool          `json:"training_enable_tensorboard"`
	TrainingTensorboardSaveImages      bool          `json:"training_tensorboard_save_images"`
	TrainingTensorboardFlushEvery      float64       `json:"training_tensorboard_flush_every"`
	SdModelCheckpoint                  string        `json:"sd_model_checkpoint"`
	SdCheckpointCache                  float64       `json:"sd_checkpoint_cache"`
	SdVaeCheckpointCache               float64       `json:"sd_vae_checkpoint_cache"`
	SdVae                              string        `json:"sd_vae"`
	SdVaeAsDefault                     bool          `json:"sd_vae_as_default"`
	InpaintingMaskWeight               float64       `json:"inpainting_mask_weight"`
	InitialNoiseMultiplier             float64       `json:"initial_noise_multiplier"`
	Img2ImgColorCorrection             bool          `json:"img2img_color_correction"`
	Img2ImgFixSteps                    bool          `json:"img2img_fix_steps"`
	Img2ImgBackgroundColor             string        `json:"img2img_background_color"`
	EnableQuantization                 bool          `json:"enable_quantization"`
	EnableEmphasis                     bool          `json:"enable_emphasis"`
	EnableBatchSeeds                   bool          `json:"enable_batch_seeds"`
	CommaPaddingBacktrack              float64       `json:"comma_padding_backtrack"`
	CLIPStopAtLastLayers               float64       `json:"CLIP_stop_at_last_layers"`
	UpcastAttn                         bool          `json:"upcast_attn"`
	UseOldEmphasisImplementation       bool          `json:"use_old_emphasis_implementation"`
	UseOldKarrasSchedulerSigmas        bool          `json:"use_old_karras_scheduler_sigmas"`
	NoDpmppSdeBatchDeterminism         bool          `json:"no_dpmpp_sde_batch_determinism"`
	UseOldHiresFixWidthHeight          bool          `json:"use_old_hires_fix_width_height"`
	InterrogateKeepModelsInMemory      bool          `json:"interrogate_keep_models_in_memory"`
	InterrogateReturnRanks             bool          `json:"interrogate_return_ranks"`
	InterrogateClipNumBeams            float64       `json:"interrogate_clip_num_beams"`
	InterrogateClipMinLength           float64       `json:"interrogate_clip_min_length"`
	InterrogateClipMaxLength           float64       `json:"interrogate_clip_max_length"`
	InterrogateClipDictLimit           float64       `json:"interrogate_clip_dict_limit"`
	InterrogateClipSkipCategories      []interface{} `json:"interrogate_clip_skip_categories"`
	InterrogateDeepbooruScoreThreshold float64       `json:"interrogate_deepbooru_score_threshold"`
	DeepbooruSortAlpha                 bool          `json:"deepbooru_sort_alpha"`
	DeepbooruUseSpaces                 bool          `json:"deepbooru_use_spaces"`
	DeepbooruEscape                    bool          `json:"deepbooru_escape"`
	DeepbooruFilterTags                string        `json:"deepbooru_filter_tags"`
	ExtraNetworksDefaultView           string        `json:"extra_networks_default_view"`
	ExtraNetworksDefaultMultiplier     float64       `json:"extra_networks_default_multiplier"`
	ExtraNetworksAddTextSeparator      string        `json:"extra_networks_add_text_separator"`
	SdHypernetwork                     string        `json:"sd_hypernetwork"`
	ReturnGrid                         bool          `json:"return_grid"`
	DoNotShowImages                    bool          `json:"do_not_show_images"`
	AddModelHashToInfo                 bool          `json:"add_model_hash_to_info"`
	AddModelNameToInfo                 bool          `json:"add_model_name_to_info"`
	DisableWeightsAutoSwap             bool          `json:"disable_weights_auto_swap"`
	SendSeed                           bool          `json:"send_seed"`
	SendSize                           bool          `json:"send_size"`
	Font                               string        `json:"font"`
	JsModalLightbox                    bool          `json:"js_modal_lightbox"`
	JsModalLightboxInitiallyZoomed     bool          `json:"js_modal_lightbox_initially_zoomed"`
	ShowProgressInTitle                bool          `json:"show_progress_in_title"`
	SamplersInDropdown                 bool          `json:"samplers_in_dropdown"`
	DimensionsAndBatchTogether         bool          `json:"dimensions_and_batch_together"`
	KeyeditPrecisionAttention          float64       `json:"keyedit_precision_attention"`
	KeyeditPrecisionExtra              float64       `json:"keyedit_precision_extra"`
	Quicksettings                      string        `json:"quicksettings"`
	HiddenTabs                         []interface{} `json:"hidden_tabs"`
	UIReorder                          string        `json:"ui_reorder"`
	UIExtraNetworksTabReorder          string        `json:"ui_extra_networks_tab_reorder"`
	Localization                       string        `json:"localization"`
	ShowProgressbar                    bool          `json:"show_progressbar"`
	LivePreviewsEnable                 bool          `json:"live_previews_enable"`
	ShowProgressGrid                   bool          `json:"show_progress_grid"`
	ShowProgressEveryNSteps            float64       `json:"show_progress_every_n_steps"`
	ShowProgressType                   string        `json:"show_progress_type"`
	LivePreviewContent                 string        `json:"live_preview_content"`
	LivePreviewRefreshPeriod           float64       `json:"live_preview_refresh_period"`
	HideSamplers                       []interface{} `json:"hide_samplers"`
	EtaDdim                            float64       `json:"eta_ddim"`
	EtaAncestral                       float64       `json:"eta_ancestral"`
	DdimDiscretize                     string        `json:"ddim_discretize"`
	SChurn                             float64       `json:"s_churn"`
	STmin                              float64       `json:"s_tmin"`
	SNoise                             float64       `json:"s_noise"`
	EtaNoiseSeedDelta                  float64       `json:"eta_noise_seed_delta"`
	AlwaysDiscardNextToLastSigma       bool          `json:"always_discard_next_to_last_sigma"`
	UniPcVariant                       string        `json:"uni_pc_variant"`
	UniPcSkipType                      string        `json:"uni_pc_skip_type"`
	UniPcOrder                         float64       `json:"uni_pc_order"`
	UniPcLowerOrderFinal               bool          `json:"uni_pc_lower_order_final"`
	PostprocessingEnableInMainUI       []interface{} `json:"postprocessing_enable_in_main_ui"`
	PostprocessingOperationOrder       []interface{} `json:"postprocessing_operation_order"`
	UpscalingMaxImagesInCache          float64       `json:"upscaling_max_images_in_cache"`
	DisabledExtensions                 []interface{} `json:"disabled_extensions"`
	SdCheckpointHash                   string        `json:"sd_checkpoint_hash"`
	SdLora                             string        `json:"sd_lora"`
	LoraApplyToOutputs                 bool          `json:"lora_apply_to_outputs"`
}

type ModelInfo struct {
	Title     string                 `json:"title"`
	ModelName string                 `json:"model_name"`
	Hash      string                 `json:"hash"`
	Sha256    string                 `json:"sha256"`
	Filename  string                 `json:"filename"`
	Config    map[string]interface{} `json:"config"`
}

type Sampler struct {
	Name    string                 `json:"name,omitempty"`
	Aliases []string               `json:"aliases,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type Upscaler struct {
	Name      string  `json:"name,omitempty"`
	ModelName string  `json:"model_name,omitempty"`
	ModelPath string  `json:"model_path,omitempty"`
	ModelURL  string  `json:"model_url,omitempty"`
	Scale     float64 `json:"scale,omitempty"`
}

type Progress struct {
	Progress     float64       `json:"progress,omitempty"`
	EtaRelative  float64       `json:"eta_relative,omitempty"`
	State        ProgressState `json:"state,omitempty"`
	CurrentImage string        `json:"current_image,omitempty"`
	Textinfo     string        `json:"textinfo,omitempty"`
}
type ProgressState struct {
	Skipped       bool   `json:"skipped,omitempty"`
	Interrupted   bool   `json:"interrupted,omitempty"`
	Job           string `json:"job,omitempty"`
	JobCount      int    `json:"job_count,omitempty"`
	JobTimestamp  string `json:"job_timestamp,omitempty"`
	JobNo         int    `json:"job_no,omitempty"`
	SamplingStep  int    `json:"sampling_step,omitempty"`
	SamplingSteps int    `json:"sampling_steps,omitempty"`
}
