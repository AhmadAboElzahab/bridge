import type { FilterGroup, FilterRule } from "@/types/api";
import type { FilterGroupNode, FilterItemNode, FilterNode } from "./types";

function itemToRule(item: FilterItemNode): FilterRule {
  return {
    id: item.id,
    type: "RULE",
    field: item.field,
    operator: item.operator,
    value: item.value,
  };
}

function groupToApi(group: FilterGroupNode): FilterGroup {
  return {
    id: group.id,
    type: "GROUP",
    conjunction: group.conjunction,
    children: group.children
      .filter((child) => {
        if (child.type === "ITEM") return child.field && child.operator;
        return child.children.length > 0;
      })
      .map((child) => nodeToApi(child)),
  };
}

function nodeToApi(node: FilterNode): FilterGroup | FilterRule {
  if (node.type === "ITEM") return itemToRule(node);
  return groupToApi(node);
}

export function toApiFilter(
  group: FilterGroupNode,
): FilterGroup | Record<string, never> {
  const filtered = group.children.filter((child) => {
    if (child.type === "ITEM") return child.field && child.operator;
    return child.children.length > 0;
  });

  if (filtered.length === 0) return {};

  return groupToApi(group);
}
