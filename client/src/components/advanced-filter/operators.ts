import type { FieldType } from "@/types/api";

export interface OperatorDef {
  label: string;
  value: string;
  noValue?: boolean;
}

export const TEXT_OPERATORS: OperatorDef[] = [
  { label: "Contains", value: "contains" },
  { label: "Does not contain", value: "notContains" },
  { label: "Is", value: "equals" },
  { label: "Is not", value: "notEquals" },
  { label: "Starts with", value: "startsWith" },
  { label: "Ends with", value: "endsWith" },
  { label: "Is empty", value: "isEmpty", noValue: true },
  { label: "Is not empty", value: "isNotEmpty", noValue: true },
];

export const NUMBER_OPERATORS: OperatorDef[] = [
  { label: "Is", value: "equals" },
  { label: "Is not", value: "notEquals" },
  { label: "Less than", value: "lessThan" },
  { label: "Less than or equal", value: "lessThanOrEqual" },
  { label: "Greater than", value: "greaterThan" },
  { label: "Greater than or equal", value: "greaterThanOrEqual" },
  { label: "Is empty", value: "isEmpty", noValue: true },
  { label: "Is not empty", value: "isNotEmpty", noValue: true },
];

export const DATE_OPERATORS: OperatorDef[] = [
  { label: "Is", value: "equals" },
  { label: "Is not", value: "notEquals" },
  { label: "Before", value: "before" },
  { label: "After", value: "after" },
  { label: "On or before", value: "onOrBefore" },
  { label: "On or after", value: "onOrAfter" },
  { label: "Is empty", value: "isEmpty", noValue: true },
  { label: "Is not empty", value: "isNotEmpty", noValue: true },
];

export const SELECT_OPERATORS: OperatorDef[] = [
  { label: "Is", value: "is" },
  { label: "Is not", value: "isNot" },
  { label: "Is any of", value: "isAnyOf" },
  { label: "Is none of", value: "isNoneOf" },
  { label: "Is empty", value: "isEmpty", noValue: true },
  { label: "Is not empty", value: "isNotEmpty", noValue: true },
];

export const BOOLEAN_OPERATORS: OperatorDef[] = [
  { label: "Is true", value: "isTrue", noValue: true },
  { label: "Is false", value: "isFalse", noValue: true },
];

export function getOperators(fieldType: FieldType): OperatorDef[] {
  switch (fieldType) {
    case "integer_field":
    case "float_field":
    case "currency_field":
    case "rating_field":
      return NUMBER_OPERATORS;
    case "date_field":
    case "datetime_field":
      return DATE_OPERATORS;
    case "boolean_field":
      return BOOLEAN_OPERATORS;
    case "single_select":
    case "radio_select":
    case "creatable_single_select":
    case "multi_select":
    case "creatable_multi_select":
    case "single_relation":
    case "multi_relation":
      return SELECT_OPERATORS;
    default:
      return TEXT_OPERATORS;
  }
}

export function getDefaultOperator(fieldType: FieldType): string {
  switch (fieldType) {
    case "integer_field":
    case "float_field":
    case "currency_field":
    case "rating_field":
      return "equals";
    case "date_field":
    case "datetime_field":
      return "equals";
    case "boolean_field":
      return "isTrue";
    case "single_select":
    case "radio_select":
    case "creatable_single_select":
    case "multi_select":
    case "creatable_multi_select":
    case "single_relation":
    case "multi_relation":
      return "is";
    default:
      return "contains";
  }
}
