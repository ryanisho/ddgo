import React from 'react';

export const StatBox = ({ label, value }) => (
    <div className="bg-gray-50 rounded-lg p-4">
        <div className="text-sm text-gray-500 mb-1">{label}</div>
        <div className="text-2xl font-semibold">{value}</div>
    </div>
);