import { Container } from '@mui/material'
import './App.css'
import { Route, Routes } from 'react-router-dom'
import HomePage from './pages/HomePage'
import Navbar from './components/Navbar'

function App() {
  return (
    <Container>
      <Navbar />

      <Routes>
        <Route path="/" element={<HomePage />} />
      </Routes>
    </Container>
  )
}

export default App
