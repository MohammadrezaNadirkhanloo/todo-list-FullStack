import { useState } from "react";

import type { PaginationState } from "@tanstack/react-table";

// import { DEFAULT_PAGE_SIZE } from "@/components/data-table/build-columns";
import { useTableQuery } from "@/hooks/useTableQuery";




export function useListTodos() {
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 20,
  });
  const [sorts, setSorts] = useState([]);
  const [search, setSearch] = useState("");

  const primarySort = sorts[0];

  const query = useTableQuery({
    endpoint: "todos",
    page: pagination.pageIndex + 1,
    perPage: pagination.pageSize,
    filter: search,
    // orderBy: primarySort?.column ,
    // orderDirection: primarySort?.direction ?? "",
  });

  return {
    Todos: query.data,
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
