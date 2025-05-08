-- 清空现有数据
DELETE FROM registrations;
DELETE FROM volunteers;
DELETE FROM activities;
DELETE FROM users;

-- 重置自增ID
DELETE FROM sqlite_sequence;

-- 创建用户数据
INSERT INTO users (id, username, password, role, status, created_at, updated_at) VALUES
(1, 'admin', '$2a$10$pazilwVr1daGU2Qu2gUpWuvrd4OCpjaAHXZm9opXrZ8s23amNIQBO', 'admin', 'active', datetime('now'), datetime('now')),
(2, 'volunteer1', '$2a$10$pazilwVr1daGU2Qu2gUpWuvrd4OCpjaAHXZm9opXrZ8s23amNIQBO', 'volunteer', 'active', datetime('now'), datetime('now')),
(3, 'volunteer2', '$2a$10$pazilwVr1daGU2Qu2gUpWuvrd4OCpjaAHXZm9opXrZ8s23amNIQBO', 'volunteer', 'active', datetime('now'), datetime('now')),
(4, 'volunteer3', '$2a$10$pazilwVr1daGU2Qu2gUpWuvrd4OCpjaAHXZm9opXrZ8s23amNIQBO', 'volunteer', 'pending', datetime('now'), datetime('now'));

-- 创建志愿者数据
INSERT INTO volunteers (id, user_id, name, phone, email, address, hours, activities, status, created_at, updated_at) VALUES
(1, 2, '张志愿', '13800138001', 'zhang@example.com', '北京市海淀区', 10, 2, 'active', datetime('now'), datetime('now')),
(2, 3, '李志愿', '13800138002', 'li@example.com', '北京市朝阳区', 5, 1, 'active', datetime('now'), datetime('now')),
(3, 4, '王志愿', '13800138003', 'wang@example.com', '北京市西城区', 0, 0, 'pending', datetime('now'), datetime('now'));

-- 创建活动数据
INSERT INTO activities (id, title, date, status, location, capacity, registered, description, created_at, updated_at) VALUES
(1, '社区清洁日', datetime('now', '+7 days'), 'upcoming', '北京市海淀区中关村街道', 20, 2, '组织社区清洁活动，美化环境', datetime('now'), datetime('now')),
(2, '敬老院慰问', datetime('now', '+14 days'), 'upcoming', '北京市朝阳区敬老院', 15, 1, '看望敬老院老人，带去节日的温暖', datetime('now'), datetime('now')),
(3, '公园植树', datetime('now', '-7 days'), 'completed', '北京市海淀区公园', 30, 30, '参与植树造林，绿化环境', datetime('now'), datetime('now'));

-- 创建报名记录数据
INSERT INTO registrations (id, activity_id, user_id, name, phone, id_card, email, emergency_contact, emergency_phone, status, create_time, updated_at) VALUES
(1, 1, 2, '张志愿', '13800138001', '110101199001011234', 'zhang@example.com', '张父', '13900139001', 'approved', datetime('now'), datetime('now')),
(2, 1, 3, '李志愿', '13800138002', '110101199001011235', 'li@example.com', '李母', '13900139002', 'approved', datetime('now'), datetime('now')),
(3, 2, 2, '张志愿', '13800138001', '110101199001011234', 'zhang@example.com', '张父', '13900139001', 'pending', datetime('now'), datetime('now'));
