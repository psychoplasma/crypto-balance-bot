import { useState } from 'react';
import { useRouter } from 'next/router';
import { login } from '../api/api';
import Layout from '../components/Layout';

const Login = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const router = useRouter();

    const handleSubmit = async (e: any) => {
        e.preventDefault();
        setError('');

        const res = await login(email, password);

        if (res.ok) {
            router.push('/user');
        } else {
            const data = await res.json();
            setError(data.message || 'Login failed');
        }
    };

    return (
        <Layout>
            <div className="login-container">
                <div className="login-card">
                    <h1>Login</h1>
                    <form onSubmit={handleSubmit}>
                        <div className="form-group">
                            <label htmlFor="email">Email:</label>
                            <input
                                type="email"
                                id="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                            />
                        </div>
                        <div className="form-group">
                            <label htmlFor="password">Password:</label>
                            <input
                                type="password"
                                id="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                            />
                        </div>
                        {error && <p className="error">{error}</p>}
                        <button type="submit" className="submit-btn">Login</button>
                    </form>
                    <p className="signup-link">
                        Don't have an account? <a href="/signup">Sign up</a>
                    </p>
                </div>
                <style jsx>{`
                    .login-container {
                        display: flex;
                        justify-content: center;
                        align-items: center;
                        min-height: calc(100vh - 200px);
                        padding: 2rem;
                    }

                    .login-card {
                        background: white;
                        padding: 2rem;
                        border-radius: 8px;
                        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
                        width: 100%;
                        max-width: 400px;
                    }

                    h1 {
                        text-align: center;
                        margin-bottom: 2rem;
                        color: #333;
                    }

                    .form-group {
                        margin-bottom: 1.5rem;
                    }

                    label {
                        display: block;
                        margin-bottom: 0.5rem;
                        color: #555;
                        font-weight: 500;
                    }

                    input {
                        width: 100%;
                        padding: 0.75rem;
                        border: 1px solid #ddd;
                        border-radius: 4px;
                        font-size: 1rem;
                        transition: border-color 0.2s ease;
                    }

                    input:focus {
                        outline: none;
                        border-color: #0070f3;
                    }

                    .error {
                        color: #e00;
                        margin-bottom: 1rem;
                        font-size: 0.875rem;
                    }

                    .submit-btn {
                        width: 100%;
                        padding: 0.75rem;
                        background-color: #0070f3;
                        color: white;
                        border: none;
                        border-radius: 4px;
                        font-size: 1rem;
                        font-weight: 500;
                        cursor: pointer;
                        transition: background-color 0.2s ease;
                    }

                    .submit-btn:hover {
                        background-color: #0060df;
                    }

                    .signup-link {
                        text-align: center;
                        margin-top: 1.5rem;
                        color: #666;
                    }

                    .signup-link a {
                        color: #0070f3;
                        text-decoration: none;
                        font-weight: 500;
                    }

                    .signup-link a:hover {
                        text-decoration: underline;
                    }
                `}</style>
            </div>
        </Layout>
    );
};

export default Login;