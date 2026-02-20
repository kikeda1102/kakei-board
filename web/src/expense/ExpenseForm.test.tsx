import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import ExpenseForm from "./ExpenseForm";
import "@testing-library/jest-dom";
import * as expenseApi from "./expenseApi";

jest.mock("./expenseApi");

const mockUseRecordExpenseMutation =
  expenseApi.useRecordExpenseMutation as jest.MockedFunction<
    typeof expenseApi.useRecordExpenseMutation
  >;

function setupMockMutation(
  overrides: { isLoading?: boolean; error?: unknown } = {},
) {
  const trigger = jest.fn().mockReturnValue({
    unwrap: () => Promise.resolve({ id: "new-id" }),
  });

  mockUseRecordExpenseMutation.mockReturnValue([
    trigger,
    {
      isLoading: overrides.isLoading ?? false,
      error: overrides.error,
      reset: jest.fn(),
    } as unknown as ReturnType<typeof expenseApi.useRecordExpenseMutation>[1],
  ]);

  return trigger;
}

test("renders form fields", () => {
  setupMockMutation();

  render(<ExpenseForm />);

  expect(screen.getByPlaceholderText("1500")).toBeInTheDocument();
  expect(screen.getByText("カテゴリ")).toBeInTheDocument();
  expect(screen.getByPlaceholderText("コンビニ")).toBeInTheDocument();
  expect(screen.getByText("記録する")).toBeInTheDocument();
});

test("submits expense with correct data", async () => {
  const user = userEvent.setup();
  const trigger = setupMockMutation();

  render(<ExpenseForm />);

  await user.type(screen.getByPlaceholderText("1500"), "2000");
  await user.type(screen.getByPlaceholderText("コンビニ"), "ランチ");
  await user.click(screen.getByText("記録する"));

  expect(trigger).toHaveBeenCalledWith(
    expect.objectContaining({
      amount: 2000,
      category: "食費",
      memo: "ランチ",
    }),
  );
});

test("shows disabled button while loading", () => {
  setupMockMutation({ isLoading: true });

  render(<ExpenseForm />);

  const button = screen.getByText("記録中...");
  expect(button).toBeDisabled();
});

test("shows error message on API failure", () => {
  setupMockMutation({ error: { status: 500, data: "fail" } });

  render(<ExpenseForm />);

  expect(screen.getByText("記録に失敗しました")).toBeInTheDocument();
});
