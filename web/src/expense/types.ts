export interface Expense {
  id: string;
  amount: number;
  category: string;
  memo: string;
  date: string;
  created_at: string;
}

export interface RecordExpenseRequest {
  amount: number;
  category: string;
  memo: string;
  date: string;
}

export interface RecordExpenseResponse {
  id: string;
}
