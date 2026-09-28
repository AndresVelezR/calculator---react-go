import { Calculator } from './calculator/Calculator'

function App() {
  return (
    <main className="app-shell">
      <header className="app-heading">
        <p className="eyebrow">A little everyday math</p>
        <h1>Calculator</h1>
        <p>From the basics to a little more.</p>
      </header>
      <Calculator />
    </main>
  )
}

export default App
