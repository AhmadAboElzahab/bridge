import * as React from "react";
import type {
  FilterGroupNode,
  FilterItemNode,
  FilterNode,
} from "@/components/advanced-filter/types";
import {
  addGroupToGroup,
  addItemToGroup,
  createEmptyRootGroup,
  createFilterGroup,
  createFilterItem,
  removeNodeFromGroup,
  reorderGroupChildren,
  updateGroupConjunction,
  updateItemInGroup,
} from "@/components/advanced-filter/utils";

export const MAX_NESTING_LEVEL = 1;

export function useAdvancedFilter(
  initialValue: FilterGroupNode | undefined,
  onChange: (ruleset: FilterGroupNode) => void,
) {
  const [ruleset, setRuleset] = React.useState<FilterGroupNode>(
    () => initialValue ?? createEmptyRootGroup(),
  );

  const onChangeRef = React.useRef(onChange);
  React.useEffect(() => { onChangeRef.current = onChange; }, [onChange]);

  const pendingChangeRef = React.useRef<FilterGroupNode | null>(null);

  const commit = React.useCallback((updater: (prev: FilterGroupNode) => FilterGroupNode) => {
    setRuleset((prev) => {
      const next = updater(prev);
      pendingChangeRef.current = next;
      return next;
    });
  }, []);

  React.useEffect(() => {
    if (pendingChangeRef.current !== null) {
      const next = pendingChangeRef.current;
      pendingChangeRef.current = null;
      onChangeRef.current(next);
    }
  });

  const addItem = React.useCallback(
    (groupId: string, field: string, operator: string) => {
      const item = createFilterItem(field, operator);
      commit((prev) => addItemToGroup(prev, groupId, item));
    },
    [commit],
  );

  const addGroup = React.useCallback(
    (parentId: string) => {
      const group = createFilterGroup();
      commit((prev) => addGroupToGroup(prev, parentId, group));
    },
    [commit],
  );

  const removeNode = React.useCallback(
    (id: string) => {
      commit((prev) => removeNodeFromGroup(prev, id));
    },
    [commit],
  );

  const updateItem = React.useCallback(
    (id: string, updates: Partial<Omit<FilterItemNode, "id" | "type">>) => {
      commit((prev) => updateItemInGroup(prev, id, updates));
    },
    [commit],
  );

  const updateConjunction = React.useCallback(
    (groupId: string, conjunction: "and" | "or") => {
      commit((prev) => updateGroupConjunction(prev, groupId, conjunction));
    },
    [commit],
  );

  const reorderChildren = React.useCallback(
    (groupId: string, newChildren: FilterNode[]) => {
      commit((prev) => reorderGroupChildren(prev, groupId, newChildren));
    },
    [commit],
  );

  const clearAll = React.useCallback(() => {
    commit(() => createEmptyRootGroup());
  }, [commit]);

  return {
    ruleset,
    addItem,
    addGroup,
    removeNode,
    updateItem,
    updateConjunction,
    reorderChildren,
    clearAll,
  };
}
