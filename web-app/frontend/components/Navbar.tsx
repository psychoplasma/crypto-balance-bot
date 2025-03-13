import React from 'react';
import Link from 'next/link';

const Navbar: React.FC = () => {
    return (
        <nav className="navbar">
            <div className="navbar-brand">
                <Link href="/">Crypto Balance Bot</Link>
            </div>
            <div className="navbar-links">
                <Link href="/signup">Sign Up</Link>
                <Link href="/login">Login</Link>
                <Link href="/user">User Dashboard</Link>
            </div>
            <style jsx>{`
                .navbar {
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    padding: 1rem 2rem;
                    background-color: #fff;
                    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
                }

                .navbar-brand {
                    font-size: 1.5rem;
                    font-weight: bold;
                }

                .navbar-brand a {
                    color: #0070f3;
                    text-decoration: none;
                }

                .navbar-links {
                    display: flex;
                    gap: 2rem;
                    align-items: center;
                }

                .navbar-links a {
                    color: #666;
                    text-decoration: none;
                    font-weight: 500;
                    transition: color 0.2s ease;
                }

                .navbar-links a:hover {
                    color: #0070f3;
                }
            `}</style>
        </nav>
    );
};

export default Navbar;