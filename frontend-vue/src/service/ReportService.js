import ApiService from './ApiService';

export const ReportService = {
    generateReport(month, year) {
        return ApiService.get(`/reports/generate?month=${month}&year=${year}`);
    }
};