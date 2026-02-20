import ExpenseForm from "./ExpenseForm";
import ExpenseList from "./ExpenseList";

export default function ExpensePage() {
  return (
    <>
      <h2>支出</h2>
      <ExpenseForm />
      <ExpenseList />
    </>
  );
}
