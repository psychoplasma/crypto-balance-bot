import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { signup } from '../../lib/api';
import Layout from '../../components/Layout';
import './signup.css';

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
            </div>
        </Layout>
    );
};

export default Signup;