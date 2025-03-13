import { useState } from 'react';
import { useRouter } from 'next/router';
import { signup } from '../api/api';
import Layout from '../components/Layout';

const Signup = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const router = useRouter();

    const handleSubmit = async (e: any) => {
        e.preventDefault();
        setError('');

        const response = await signup(email, password);

        if (response.ok) {
            router.push('/login');
        } else {
            const data = await response.json();
            setError(data.message || 'Something went wrong');
        }
    };

    return (
        <Layout>
            <div className="signup-container">
                <div className="signup-card">
                    <h1>Sign Up</h1>
                    <form onSubmit={handleSubmit}>
                        <div className="form-group">
                            <label htmlFor="email">Email</label>
                            <input
                                type="email"
                                id="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                            />
                        </div>
                        <div className="form-group">
                            <label htmlFor="password">Password</label>
                            <input
                                type="password"
                                id="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                            />
                        </div>
                        {error && <p className="error">{error}</p>}
                        <button type="submit" className="submit-btn">Sign Up</button>
                    </form>
                    <p className="login-link">
                        Already have an account? <a href="/login">Log in</a>
                    </p>
                </div>
                <style jsx>{`
                    .signup-container {
                        display: flex;
                        justify-content: center;
                        align-items: center;
                        min-height: calc(100vh - 200px);
                        padding: 2rem;
                    }

                    .signup-card {
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

                    .login-link {
                        text-align: center;
                        margin-top: 1.5rem;
                        color: #666;
                    }

                    .login-link a {
                        color: #0070f3;
                        text-decoration: none;
                        font-weight: 500;
                    }

                    .login-link a:hover {
                        text-decoration: underline;
                    }
                `}</style>
            </div>
        </Layout>
    );
};

export default Signup;