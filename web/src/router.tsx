import { createBrowserRouter } from "react-router-dom";
import Layout from "./shared/components/Layout";
import ExpensePage from "./expense/ExpensePage";
import SummaryPage from "./pages/SummaryPage";
import BudgetPage from "./pages/BudgetPage";
import ScorePage from "./pages/ScorePage";

export const router = createBrowserRouter([
  {
    element: <Layout />,
    children: [
      { index: true, element: <ExpensePage /> },
      { path: "summary", element: <SummaryPage /> },
      { path: "budget", element: <BudgetPage /> },
      { path: "score", element: <ScorePage /> },
    ],
  },
]);
