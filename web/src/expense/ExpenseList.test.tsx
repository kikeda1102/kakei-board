import { render, screen } from "@testing-library/react";
import ExpenseList from "./ExpenseList";
import "@testing-library/jest-dom";
import * as expenseApi from "./expenseApi";

jest.mock("./expenseApi");

const mockUseListExpensesQuery =
  expenseApi.useListExpensesQuery as jest.MockedFunction<
    typeof expenseApi.useListExpensesQuery
  >;

test("displays expenses returned by the API", () => {
  mockUseListExpensesQuery.mockReturnValue({
    data: [
      {
        id: "abc",
        amount: 1500,
        category: "食費",
        memo: "コンビニ",
        date: "2026-02-20",
        created_at: "2026-02-20T00:00:00Z",
      },
    ],
    isLoading: false,
    error: undefined,
    refetch: jest.fn(),
  } as unknown as ReturnType<typeof expenseApi.useListExpensesQuery>);

  render(<ExpenseList />);

  expect(screen.getByText("食費")).toBeInTheDocument();
  expect(screen.getByText("コンビニ")).toBeInTheDocument();
  expect(screen.getByText("1,500円")).toBeInTheDocument();
  expect(screen.getByText("2026-02-20")).toBeInTheDocument();
});

test("displays empty message when no expenses", () => {
  mockUseListExpensesQuery.mockReturnValue({
    data: [],
    isLoading: false,
    error: undefined,
    refetch: jest.fn(),
  } as unknown as ReturnType<typeof expenseApi.useListExpensesQuery>);

  render(<ExpenseList />);

  expect(screen.getByText("支出の記録がありません")).toBeInTheDocument();
});

test("displays loading message", () => {
  mockUseListExpensesQuery.mockReturnValue({
    data: undefined,
    isLoading: true,
    error: undefined,
    refetch: jest.fn(),
  } as unknown as ReturnType<typeof expenseApi.useListExpensesQuery>);

  render(<ExpenseList />);

  expect(screen.getByText("読み込み中...")).toBeInTheDocument();
});

test("displays error message on API failure", () => {
  mockUseListExpensesQuery.mockReturnValue({
    data: undefined,
    isLoading: false,
    error: { status: 500, data: "fail" },
    refetch: jest.fn(),
  } as unknown as ReturnType<typeof expenseApi.useListExpensesQuery>);

  render(<ExpenseList />);

  expect(screen.getByText("一覧の取得に失敗しました")).toBeInTheDocument();
});
