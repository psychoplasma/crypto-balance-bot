import React from 'react';
import Link from 'next/link';
import './Navbar.css';

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
        </nav>
    );
};

export default Navbar;
