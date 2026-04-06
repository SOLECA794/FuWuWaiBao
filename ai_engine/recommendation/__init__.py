<<<<<<< HEAD
﻿"""Recommendation service package for AI teaching platform."""
=======
"""Recommendation service package for AI teaching platform."""
>>>>>>> d17b116d297b507f8a5227ba4474640a7e13e8e0

from .adapters import BaseSourceAdapter, LocalDBAdapter, WebVideoAdapter
from .schemas import RequestDTO, ResourceDTO, TeacherFeedbackDTO
from .service import RecommendationEngine, create_recommendation_app

__all__ = [
	"BaseSourceAdapter",
	"LocalDBAdapter",
	"WebVideoAdapter",
	"RequestDTO",
	"ResourceDTO",
	"TeacherFeedbackDTO",
	"RecommendationEngine",
	"create_recommendation_app",
]
