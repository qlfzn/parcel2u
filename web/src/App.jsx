import { useState } from "react"

function App() {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [role, setRole] = useState("student")
  const [isSignUp, setIsSignUp] = useState(false)

  async function handleLogin(e) {
    e.preventDefault()
    try {
      const response = await fetch("http://localhost:8080/auth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });
      if (!response.ok) {
        throw new Error("Login failed");
      }
      const data = await response.json();
      alert("Login successful! Welcome, " + data.user.username);
    } catch (error) {
      alert("Login failed: " + error.message);
    }
  }

  async function handleSignUp(e) {
    e.preventDefault()
    try {
      const response = await fetch("http://localhost:8080/auth/users", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password, role }),
      });
      if (!response.ok) {
        throw new Error("Sign up failed");
      }
      const data = await response.json();
      alert("Sign up successful! Welcome, " + data.user.username);
      setIsSignUp(false);  
    } catch (error) {
      alert("Sign up failed: " + error.message);
    }
  }

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="flex">
        <form
          className="bg-white p-8 rounded shadow-md flex flex-col gap-4"
          onSubmit={isSignUp ? handleSignUp : handleLogin}
        >
          <h2 className="text-2xl font-bold mb-4">{isSignUp ? "Sign Up" : "Login"} form</h2>
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={e => setUsername(e.target.value)}
            className="border p-2 rounded"
            required
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            className="border p-2 rounded"
            required
          />
          {isSignUp && (
            <select
              value={role}
              onChange={e => setRole(e.target.value)}
              className="bg-gray-100 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5"
              required
            >
              <option value="student">Student</option>
              <option value="dispatch">Dispatch</option>
            </select>
          )}
          <button
            type="submit"
            className="bg-black text-white p-2 rounded hover:bg-gray-500"
          >
            {isSignUp ? "Sign Up" : "Login"}
          </button>
          <button
            type="button"
            className="underline text-blue-600 mt-2"
            onClick={() => setIsSignUp(!isSignUp)}
          >
            {isSignUp ? "Already have an account? Login" : "Don't have an account? Sign Up"}
          </button>
        </form>
      </div>
    </div>
  )
}

export default App
