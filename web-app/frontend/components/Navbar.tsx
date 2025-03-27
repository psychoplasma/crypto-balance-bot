'use client';

import React from 'react';
import Link from 'next/link';
import { useAuth } from '@/context/AuthContext';
import { useRouter } from 'next/navigation';
import './Navbar.css';

const Navbar: React.FC = () => {
  const { isAuthenticated, authLogout } = useAuth();
  const router = useRouter();

  const handleLogout = async () => {
    await authLogout();
    router.push('/login');
  };

  return (
    <nav className="navbar">
      <div className="navbar-brand">
        {isAuthenticated ? (
          <Link href="/user">Crypto Balance Bot</Link>
        ) : (
        <Link href="/">Crypto Balance Bot</Link>
        )}
      </div>
      <div className="navbar-links">
        {isAuthenticated ? (
          <>
            <Link href="/user">User Dashboard</Link>
            <button onClick={handleLogout} className="nav-button">Logout</button>
          </>
        ) : (
          <>
            <Link href="/signup">Signup</Link>
            <Link href="/login">Login</Link>
          </>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
