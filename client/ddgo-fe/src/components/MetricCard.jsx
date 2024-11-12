import React from 'react';

export const MetricCard = ({ title, children, className = '' }) => (
    <div className={`metric-card p-6 ${className}`}>
        <h2 className="text-xl font-semibold mb-4 text-gray-800">{title}</h2>
        {children}
    </div>
);