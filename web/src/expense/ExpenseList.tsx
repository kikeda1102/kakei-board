import styled from "styled-components";
import { useListExpensesQuery } from "./expenseApi";

const Table = styled.table`
  width: 100%;
  border-collapse: collapse;
  font-size: 0.95rem;
`;

const Th = styled.th`
  text-align: left;
  padding: 0.5rem 0.75rem;
  border-bottom: 2px solid #333;
  white-space: nowrap;
`;

const Td = styled.td`
  padding: 0.5rem 0.75rem;
  border-bottom: 1px solid #2a2a2a;
`;

const AmountCell = styled(Td)`
  text-align: right;
  font-variant-numeric: tabular-nums;
`;

const Message = styled.p`
  color: #888;
  font-size: 0.9rem;
`;

export default function ExpenseList() {
  const { data: expenses, isLoading, error } = useListExpensesQuery();

  if (isLoading) return <Message>読み込み中...</Message>;
  if (error) return <Message>一覧の取得に失敗しました</Message>;
  if (!expenses || expenses.length === 0)
    return <Message>支出の記録がありません</Message>;

  return (
    <Table>
      <thead>
        <tr>
          <Th>日付</Th>
          <Th>カテゴリ</Th>
          <Th>メモ</Th>
          <Th style={{ textAlign: "right" }}>金額</Th>
        </tr>
      </thead>
      <tbody>
        {expenses.map((e) => (
          <tr key={e.id}>
            <Td>{e.date}</Td>
            <Td>{e.category}</Td>
            <Td>{e.memo}</Td>
            <AmountCell>{e.amount.toLocaleString()}円</AmountCell>
          </tr>
        ))}
      </tbody>
    </Table>
  );
}
