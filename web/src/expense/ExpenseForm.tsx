import { type FormEvent, useState } from "react";
import styled from "styled-components";
import { useRecordExpenseMutation } from "./expenseApi";

const CATEGORY_PRESETS = [
  "食費",
  "交通費",
  "住居費",
  "光熱費",
  "娯楽費",
  "その他",
] as const;

const TODAY = new Date().toISOString().slice(0, 10);

const Form = styled.form`
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  align-items: flex-end;
  margin-bottom: 1.5rem;
`;

const Field = styled.label`
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  font-size: 0.85rem;
`;

const Input = styled.input`
  padding: 0.5rem;
  border: 1px solid #444;
  border-radius: 4px;
  background: #1a1a1a;
  color: inherit;
  font-size: 0.95rem;
`;

const Select = styled.select`
  padding: 0.5rem;
  border: 1px solid #444;
  border-radius: 4px;
  background: #1a1a1a;
  color: inherit;
  font-size: 0.95rem;
`;

const SubmitButton = styled.button`
  padding: 0.5rem 1.5rem;
  background: #0f3460;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.95rem;

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
`;

const ErrorMessage = styled.p`
  color: #ff6b6b;
  font-size: 0.85rem;
  width: 100%;
  margin: 0;
`;

export default function ExpenseForm() {
  const [amount, setAmount] = useState("");
  const [category, setCategory] = useState(CATEGORY_PRESETS[0]);
  const [memo, setMemo] = useState("");
  const [date, setDate] = useState(TODAY);

  const [recordExpense, { isLoading, error }] = useRecordExpenseMutation();

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();

    const parsed = Number(amount);
    if (!Number.isInteger(parsed) || parsed <= 0) return;

    await recordExpense({
      amount: parsed,
      category,
      memo,
      date,
    }).unwrap();

    setAmount("");
    setMemo("");
  };

  return (
    <Form onSubmit={handleSubmit}>
      <Field>
        金額
        <Input
          type="number"
          min="1"
          step="1"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
          required
          placeholder="1500"
        />
      </Field>

      <Field>
        カテゴリ
        <Select value={category} onChange={(e) => setCategory(e.target.value)}>
          {CATEGORY_PRESETS.map((c) => (
            <option key={c} value={c}>
              {c}
            </option>
          ))}
        </Select>
      </Field>

      <Field>
        メモ
        <Input
          type="text"
          value={memo}
          onChange={(e) => setMemo(e.target.value)}
          placeholder="コンビニ"
        />
      </Field>

      <Field>
        日付
        <Input
          type="date"
          value={date}
          onChange={(e) => setDate(e.target.value)}
          required
        />
      </Field>

      <SubmitButton type="submit" disabled={isLoading}>
        {isLoading ? "記録中..." : "記録する"}
      </SubmitButton>

      {error && <ErrorMessage>記録に失敗しました</ErrorMessage>}
    </Form>
  );
}
