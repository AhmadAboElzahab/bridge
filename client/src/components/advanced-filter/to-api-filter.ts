import type { FilterGroup, FilterRule } from "@/types/api";
import type { FilterGroupNode, FilterItemNode, FilterNode } from "./types";

// Converts a saved API FilterGroup back into the UI FilterGroupNode tree.
export function fromApiFilter(group: FilterGroup): FilterGroupNode {
  return {
    id: group.id,
    type: "GROUP",
    conjunction: group.conjunction,
    children: group.children.map((child): FilterItemNode | FilterGroupNode => {
      if (child.type === "GROUP") return fromApiFilter(child);
      return {
        id: child.id,
        type: "ITEM",
        field: child.field,
        operator: child.operator,
        value: child.value,
      };
    }),
  };
}

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
