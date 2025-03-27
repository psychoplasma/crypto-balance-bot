import Layout from '@/components/Layout';
import LoginForm from '@/components/LoginForm';
import PublicRoute from '@/components/PublicRoute';
import './login.css';

export default function Login() {
  return (
    <PublicRoute>
      <Layout>
        <div className="login-container">
          <div className="login-card">
            <h1>Login</h1>
            <LoginForm />
            <p className="signup-link">
              Don't have an account? <a href="/signup">Sign up</a>
            </p>
          </div>
        </div>
      </Layout>
    </PublicRoute>
  );
};
