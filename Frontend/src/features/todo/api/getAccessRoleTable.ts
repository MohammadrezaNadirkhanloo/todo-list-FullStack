import { useState } from "react";

import type { PaginationState } from "@tanstack/react-table";

import type { SortItem } from "@/components/data-table";
import { useTableQuery } from "@/hooks/useTableQuery";
import type { AccessRole } from "../types";
import { DEFAULT_PAGE_SIZE } from "@/components/data-table/build-columns";




export function useAccessRoleTable() {
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: DEFAULT_PAGE_SIZE,
  });
  const [sorts, setSorts] = useState<SortItem[]>([]);
  const [search, setSearch] = useState("");

  const primarySort = sorts[0];

  const query = useTableQuery<AccessRole>({
    endpoint: "/acl/roles",
    page: pagination.pageIndex + 1,
    perPage: pagination.pageSize,
    filter: search,
    orderBy: primarySort?.column ,
    orderDirection: primarySort?.direction ?? "",
  });

  return {
    AccessRole: query.data,
    rowCount: query.total,
    pageCount: query.totalPages,
    isLoading: query.isLoading,
    isFetching: query.isFetching,
    isError: query.isError,
    error: query.error,
    refetch: query.refetch,

    pagination,
    onPaginationChange: setPagination,
    sorts,
    onSortsChange: setSorts,
    onGlobalSearchChange: setSearch,
  };
}
