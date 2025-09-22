# Dynamic Evaluation System Usage Guide

## Overview

The evaluation system has been enhanced to support dynamic question and answer generation. Each evaluation attempt now generates a unique set of questions and answer options from the question bank, making assessments more robust and preventing memorization.

## How It Works

### 1. Evaluation Configuration

When creating an evaluation, you can now configure:
- **Question Count**: Number of questions to show per attempt (selected randomly from question bank)
- **Answer Options Count**: Number of answer options to show per question (selected randomly from answer bank)
- **Other settings**: Passing score, time limit, max attempts, etc.

### 2. Question Bank Setup

- Create questions for the evaluation as before
- Each question should have a large pool of answers (both correct and incorrect)
- The system will randomly select from this pool for each attempt

### 3. Answer Distribution Rules

For each generated question:
- **Minimum 1 correct answer**: At least one correct option will always be included
- **Maximum 50% correct**: At most half of the answer options will be correct
- **Random selection**: Both questions and answers are randomly selected for each attempt

### 4. Attempt Generation Process

When a user starts an evaluation attempt:

1. **Question Selection**: System randomly selects X questions from the question bank
2. **Answer Generation**: For each question, randomly selects Y answer options
3. **Distribution Validation**: Ensures proper correct/incorrect answer ratio
4. **Attempt Creation**: Stores the generated content for this specific attempt

### 5. Independent Scoring

Each attempt is scored based on its generated content:
- No dependency on original question/answer structure
- Self-contained validation logic
- Consistent scoring across different generated variations

## API Changes

### Backend Models

```go
// Enhanced Evaluation model
type Evaluation struct {
    // ... existing fields
    QuestionCount      int `json:"question_count" gorm:"not null"`
    AnswerOptionsCount int `json:"answer_options_count" gorm:"not null;default:4"`
    // ... other fields
}

// New AttemptQuestion structure
type AttemptQuestion struct {
    ID            uint                  `json:"id"`
    Text          string                `json:"text"`
    Type          enums.QuestionType    `json:"type"`
    Points        int                   `json:"points"`
    OriginalID    uint                  `json:"original_id"`
    AnswerOptions []AttemptAnswerOption `json:"answer_options"`
}

// Updated EvaluationAttempt
type EvaluationAttempt struct {
    // ... existing fields  
    Questions AttemptQuestions `json:"questions" gorm:"type:json"`
    Answers   AttemptAnswers   `json:"answers" gorm:"type:json"`
    // ... other fields
}
```

### Frontend Types

```typescript
// Enhanced Evaluation interface
export interface Evaluation extends BaseEntity {
    // ... existing fields
    question_count: number;
    answer_options_count: number;
    // ... other fields
}

// New AttemptQuestion interface
export interface AttemptQuestion {
    id: number;
    text: string;
    type: QuestionType;
    points: number;
    original_id: number;
    answer_options: AttemptAnswerOption[];
}
```

## Usage Examples

### Creating a Dynamic Evaluation

```json
{
    "title": "Dynamic Math Quiz",
    "description": "Randomized math questions",
    "question_count": 10,
    "answer_options_count": 4,
    "passing_score": 70,
    "max_attempts": 3,
    "time_limit": 30,
    "module_id": 123
}
```

### Generated Attempt Structure

```json
{
    "id": 456,
    "user_id": 789,
    "evaluation_id": 123,
    "questions": [
        {
            "id": 1,
            "text": "What is 2 + 2?",
            "type": "single_choice",
            "points": 5,
            "original_id": 42,
            "answer_options": [
                {"id": 1, "text": "3", "is_correct": false},
                {"id": 2, "text": "4", "is_correct": true},
                {"id": 3, "text": "5", "is_correct": false},
                {"id": 4, "text": "6", "is_correct": false}
            ]
        }
    ],
    "answers": [],
    "total_points": 50,
    "started_at": "2024-01-15T10:00:00Z"
}
```

## Best Practices

1. **Question Pool Size**: Create at least 2x more questions than the configured question_count for good randomization
2. **Answer Pool Size**: Create at least 2x more answers per question than the configured answer_options_count
3. **Balanced Difficulty**: Ensure your question pool has balanced difficulty levels
4. **Clear Instructions**: Make sure questions are clear and self-contained
5. **Regular Review**: Review and update question banks periodically

## Migration Notes

- Existing evaluations will work with the new system (default answer_options_count = 4)
- Old attempt data structure is automatically handled
- No data migration required - changes are backward compatible
- Frontend forms now include the new answer options count field

## Benefits

- **Increased Security**: Reduces cheating through question memorization
- **Fair Assessment**: Each attempt is unique while maintaining consistent difficulty
- **Scalability**: Easy to expand question banks without affecting existing attempts
- **Flexibility**: Configurable question and answer counts per evaluation
- **Reliability**: Self-contained scoring independent of question bank changes