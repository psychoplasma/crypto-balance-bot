import React from 'react';
import Link from 'next/link';
import Layout from '@/components/Layout';
import PublicRoute from '@/components/PublicRoute';
import './page.css';

const Home: React.FC = () => {
  return (
    <PublicRoute>
      <Layout>
        <div className="container">
          <h1>Welcome to Crypto Balance Bot</h1>
          <p>Your subscription system for monitoring blockchain account movements.</p>
          <div className="actions">
            <Link href="/signup" className="btn">
              Sign Up
            </Link>
            <Link href="/login" className="btn">
              Login
            </Link>
          </div>
        </div>
      </Layout>
    </PublicRoute>
  );
};

export default Home;
