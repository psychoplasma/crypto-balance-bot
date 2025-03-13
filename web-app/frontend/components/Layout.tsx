import React from 'react';
import Navbar from './Navbar';

const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    return (
        <div>
            <Navbar />
            <main className="main-content">{children}</main>
            <footer className="footer">
                <p>&copy; {new Date().getFullYear()} Crypto Balance Bot. All rights reserved.</p>
            </footer>

            <style jsx global>{`
                .layout {
                    display: flex;
                    flex-direction: column;
                    min-height: 100vh;
                }

                .main-content {
                    flex: 1;
                }

                .footer {
                    padding: 1.5rem;
                    text-align: center;
                    background-color: #f7f7f7;
                    border-top: 1px solid #eaeaea;
                    color: #666;
                }

                .footer p {
                    margin: 0;
                    font-size: 0.9rem;
                }
            `}</style>
        </div>
    );
};

export default Layout;