import React from 'react';
import Link from 'next/link';
import Layout from '../components/Layout';

const Home: React.FC = () => {
    return (
        <Layout>
            <div className="container">
                <h1>Welcome to Crypto Balance Bot</h1>
                <p>Your subscription system for monitoring blockchain account movements.</p>
                <div className="actions">
                    <Link href="/signup">
                        <a className="btn">Sign Up</a>
                    </Link>
                    <Link href="/login">
                        <a className="btn">Login</a>
                    </Link>
                </div>
            </div>
            <style jsx global>{`
                .container {
                    text-align: center;
                    margin-top: 50px;
                }
                .actions {
                    margin-top: 20px;
                }
                .btn {
                    margin: 0 10px;
                    padding: 10px 20px;
                    background-color: #0070f3;
                    color: white;
                    border: none;
                    border-radius: 5px;
                    text-decoration: none;
                    font-weight: bold;
                }
                .btn:hover {
                    background-color: #005bb5;
                }
            `}</style>
        </Layout>
    );
};

export default Home;