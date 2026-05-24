import { generateId } from "@/lib/id";
import type { FilterGroupNode, FilterItemNode, FilterNode } from "./types";

export function createFilterItem(
  field: string,
  operator: string,
): FilterItemNode {
  return { id: generateId(), type: "ITEM", field, operator, value: "" };
}

export function createFilterGroup(): FilterGroupNode {
  return { id: generateId(), type: "GROUP", conjunction: "and", children: [] };
}

export function createEmptyRootGroup(): FilterGroupNode {
  return createFilterGroup();
}

export function findNodeById(
  group: FilterGroupNode,
  id: string,
): FilterNode | null {
  if (group.id === id) return group;
  for (const child of group.children) {
    if (child.id === id) return child;
    if (child.type === "GROUP") {
      const found = findNodeById(child, id);
      if (found) return found;
    }
  }
  return null;
}

export function updateItemInGroup(
  group: FilterGroupNode,
  id: string,
  updates: Partial<Omit<FilterItemNode, "id" | "type">>,
): FilterGroupNode {
  return {
    ...group,
    children: group.children.map((child) => {
      if (child.id === id && child.type === "ITEM") {
        return { ...child, ...updates };
      }
      if (child.type === "GROUP") {
        return updateItemInGroup(child, id, updates);
      }
      return child;
    }),
  };
}

export function removeNodeFromGroup(
  group: FilterGroupNode,
  id: string,
): FilterGroupNode {
  return {
    ...group,
    children: group.children
      .filter((child) => child.id !== id)
      .map((child) =>
        child.type === "GROUP" ? removeNodeFromGroup(child, id) : child,
      ),
  };
}

export function addItemToGroup(
  group: FilterGroupNode,
  groupId: string,
  item: FilterItemNode,
): FilterGroupNode {
  if (group.id === groupId) {
    return { ...group, children: [...group.children, item] };
  }
  return {
    ...group,
    children: group.children.map((child) =>
      child.type === "GROUP" ? addItemToGroup(child, groupId, item) : child,
    ),
  };
}

export function addGroupToGroup(
  group: FilterGroupNode,
  parentId: string,
  newGroup: FilterGroupNode,
): FilterGroupNode {
  if (group.id === parentId) {
    return { ...group, children: [...group.children, newGroup] };
  }
  return {
    ...group,
    children: group.children.map((child) =>
      child.type === "GROUP"
        ? addGroupToGroup(child, parentId, newGroup)
        : child,
    ),
  };
}

export function updateGroupConjunction(
  group: FilterGroupNode,
  groupId: string,
  conjunction: "and" | "or",
): FilterGroupNode {
  if (group.id === groupId) return { ...group, conjunction };
  return {
    ...group,
    children: group.children.map((child) =>
      child.type === "GROUP"
        ? updateGroupConjunction(child, groupId, conjunction)
        : child,
    ),
  };
}

export function reorderGroupChildren(
  group: FilterGroupNode,
  groupId: string,
  newChildren: FilterNode[],
): FilterGroupNode {
  if (group.id === groupId) {
    return { ...group, children: newChildren };
  }
  return {
    ...group,
    children: group.children.map((child) =>
      child.type === "GROUP"
        ? reorderGroupChildren(child, groupId, newChildren)
        : child,
    ),
  };
}

export function countActiveFilters(group: FilterGroupNode): number {
  let count = 0;
  for (const child of group.children) {
    if (child.type === "ITEM") {
      if (child.field && child.operator) count++;
    } else {
      count += countActiveFilters(child);
    }
  }
  return count;
}
