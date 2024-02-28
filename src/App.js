import {
  BrowserRouter,
  Routes,
  Route,
} from 'react-router-dom'
import styled from 'styled-components'
import FourOhFour from './404'
import Login from './Login'
import Main from './Main'

const PageWrapper = styled('div')`
  padding: 24px 70px;
`


function Router() {

  return (
    <Routes>
      <Route index element={<Main />} />
      <Route path="login" element={<Login />} />
      <Route path="*" element={<FourOhFour />} />
    </Routes>
  )
}

export default function App() {
  return (
      <BrowserRouter>
            <PageWrapper>
              <Router />
            </PageWrapper>
      </BrowserRouter>
  )
}
