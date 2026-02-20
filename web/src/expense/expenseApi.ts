import { baseApi } from "../store/baseApi";
import type {
  Expense,
  RecordExpenseRequest,
  RecordExpenseResponse,
} from "./types";

interface ListExpensesParams {
  limit?: number;
  offset?: number;
}

export const expenseApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    listExpenses: builder.query<Expense[], ListExpensesParams | void>({
      query: (params) => ({
        url: "/expenses",
        params: params ?? undefined,
      }),
      providesTags: ["Expense"],
    }),
    recordExpense: builder.mutation<RecordExpenseResponse, RecordExpenseRequest>(
      {
        query: (body) => ({
          url: "/expenses",
          method: "POST",
          body,
        }),
        invalidatesTags: ["Expense"],
      },
    ),
  }),
});

export const { useListExpensesQuery, useRecordExpenseMutation } = expenseApi;
