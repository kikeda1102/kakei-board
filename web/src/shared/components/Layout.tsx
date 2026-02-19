import { NavLink, Outlet } from "react-router-dom";
import styled from "styled-components";

const Nav = styled.nav`
  display: flex;
  gap: 0;
  background: #1a1a2e;
  padding: 0;
`;

const StyledNavLink = styled(NavLink)`
  color: #ffffffcc;
  text-decoration: none;
  padding: 0.75rem 1.25rem;
  font-size: 0.95rem;
  transition: background 0.2s;

  &:hover {
    background: #16213e;
  }

  &.active {
    color: #ffffff;
    background: #0f3460;
    font-weight: 600;
  }
`;

const Main = styled.main`
  padding: 1.5rem;
`;

export default function Layout() {
  return (
    <>
      <Nav>
        <StyledNavLink to="/" end>
          支出
        </StyledNavLink>
        <StyledNavLink to="/summary">サマリー</StyledNavLink>
        <StyledNavLink to="/budget">予算</StyledNavLink>
        <StyledNavLink to="/score">スコア</StyledNavLink>
      </Nav>
      <Main>
        <Outlet />
      </Main>
    </>
  );
}
