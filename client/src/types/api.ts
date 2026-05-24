// ─── Auth ─────────────────────────────────────────────────────────────────────

export interface SigninRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
}

// ─── Common ───────────────────────────────────────────────────────────────────

export interface Country {
  value: number;
  label: string;
  code: string;
}

export interface Language {
  value: number;
  label: string;
  iso639_3: string;
}

export interface Skill {
  value: number;
  label: string;
}

export interface ErrorResponse {
  error: string;
  code?: string;
  details?: string;
}

// ─── User ─────────────────────────────────────────────────────────────────────

export interface User {
  id: number;
  first_name: string;
  last_name: string;
  email: string;
  avatar: string;
  blurhash: string;
  date_of_birth: string;
  created_at: string;
  updated_at: string;
}

export interface UserInput {
  first_name: string;
  email: string;
}

// ─── Maid ─────────────────────────────────────────────────────────────────────

export interface Maid {
  id: number;
  first_name: string;
  last_name: string;
  email: string;
  phone_number: string;
  whatsapp_number: string;
  date_of_birth: string;
  gender: string;
  religion: string;
  nationality_id: number;
  nationality: Country;
  bio: string;
  profile_picture: string;
  profile_picture_hash: string;
  video_introduction: string;
  current_location: string;
  preferred_work_location: string;
  work_preferences: string;
  availability: string;
  years_of_experience: number;
  expected_salary: number;
  skills: Skill[];
  languages: Language[];
  visa_status: string;
  status: string;
  is_profile_visible: boolean;
  is_email_verified: boolean;
  is_phone_verified: boolean;
  subscription_plan: string;
  subscription_status: string;
  created_at: string;
  updated_at: string;
}

export interface MaidInput {
  first_name: string;
  email: string;
}

// ─── Tabs ─────────────────────────────────────────────────────────────────────

export type FieldType =
  | "string_field"
  | "integer_field"
  | "float_field"
  | "currency_field"
  | "rating_field"
  | "boolean_field"
  | "date_field"
  | "datetime_field"
  | "single_select"
  | "radio_select"
  | "creatable_single_select"
  | "multi_select"
  | "creatable_multi_select"
  | "single_relation"
  | "multi_relation"
  | "rich_field"
  | "array_field"
  | "url_field"
  | "social_url_field"
  | "link_field"
  | "phone_field"
  | "email_field"
  | "image_field"
  | "file_field"
  | "attachments_field"
  | "password_field"
  | "formula_field"
  | "video_field";

export interface FieldOption {
  label: string;
  value: number | string;
}

export interface FormField {
  id: number;
  label: string;
  field_key: string;
  model_name: string;
  form_field_type: FieldType;
  form_help_text: string;
  form_is_required: boolean;
  form_order: number;
  form_stage: string;
  form_width: number;
  data_source: string;
  table_is_visible: boolean;
  table_is_pinned: boolean;
  table_order: number;
  options?: FieldOption[];
}

export interface UserTabColumn {
  field_key: string;
  form_field_id: number;
  locked: boolean;
  order: number;
  visible: boolean;
  width: number;
}

export interface UserTab {
  id: number;
  model_name: string;
  tab_name: string;
  is_default: boolean;
  search_term: string;
  filters: FilterGroup | Record<string, never>;
  columns: UserTabColumn[];
}

export interface TabsResponse {
  form_fields: FormField[];
  tabs: UserTab[];
}

export interface CreateTabInput {
  model_name: string;
  tab_name: string;
}

export interface UpdateTabInput {
  tab_name: string;
}

// ─── Filters ──────────────────────────────────────────────────────────────────

export interface FilterRule {
  id: string;
  type: "RULE";
  field: string;
  operator: string;
  value: unknown;
}

export interface FilterGroup {
  id: string;
  type: "GROUP";
  conjunction: "and" | "or";
  children: (FilterRule | FilterGroup)[];
}

// ─── Index Payload ────────────────────────────────────────────────────────────

export interface IndexPayload {
  tab_id: number;
  model?: string;
  filters: FilterGroup | Record<string, never>;
  search_term: string;
  columns: UserTabColumn[];
  page: number;
  size: number;
  search: string;
}

export interface IndexResponse<T = Record<string, unknown>> {
  data: T[];
  meta: {
    totalRowCount: number;
  };
}
